package service

import (
	"context"
	"delayed_notifier/internal/models"
	"delayed_notifier/internal/ports"
	"delayed_notifier/pkg/dlq"
	"delayed_notifier/pkg/types"
	"fmt"
	"github.com/wb-go/wbf/zlog"
	"golang.org/x/sync/errgroup"
	"time"
)

// SenderService performs the background operations of sending notifications to workers when datetime.now() is close to their publication_at
//
//	  go service.Run(ctx)
//	  ... // somewhere
//	  if notification.PublicationAt < service.WhenNextFetch() {
//		   service.QuickSend(ctx, notification)
//	  }
type SenderService struct {
	fetchPeriod      time.Duration
	fetchMaxDiapason time.Duration

	// publisherRepo publishes a notification or a batch into MQ
	publisherRepo ports.NotificationPublisherRepository

	// storageFetcherRepo fetches a batch of to-send notifications on request, service should query it in advance
	//
	// returns datetime for next fetch to be performed
	storageFetcherRepo ports.NotificationFetcherRepository

	// nextFetchIsAt is just for other services to know when will we wake up and gather notifications again
	nextFetchIsAt time.Time
}

// NewSenderService creates a new SenderService
func NewSenderService(fetchPeriod time.Duration, fetchMaxDiapason time.Duration,
	publisher ports.NotificationPublisherRepository, fetcher ports.NotificationFetcherRepository) *SenderService {
	return &SenderService{
		fetchPeriod:        fetchPeriod,
		fetchMaxDiapason:   fetchMaxDiapason,
		publisherRepo:      publisher,
		storageFetcherRepo: fetcher,
	}
}

// Run is the main blocking method. It tracks
func (s *SenderService) Run(ctx context.Context) {
	ticker := time.NewTicker(s.fetchPeriod)
	defer ticker.Stop()

out:
	for {
		select {
		case <-ctx.Done():
			break out
		case <-ticker.C:

			// life cycle
			now := time.Now()
			s.nextFetchIsAt = now.Add(s.fetchPeriod)

			// step 1. Get batch
			batch, err := s.storageFetcherRepo.Fetch(ctx, types.NewDateTime(now.Add(s.fetchMaxDiapason)))
			if err != nil {
				zlog.Logger.Error().Err(fmt.Errorf("failed to fetch batch for sending: %w", err)).Msg("error in SenderService loop")
				continue
			}

			// step 2. Send it
			err = s.SendBatch(ctx, batch)
			if err != nil {
				zlog.Logger.Error().Err(fmt.Errorf("failed to send batch: %w", err)).Msg("error in SenderService loop")
				continue
			}
		}
	}
}

// QuickSend sends 1 notification
//
// should be called when an ASAP notification is created (it's publication datetime is too close to “now()“)
func (s *SenderService) QuickSend(ctx context.Context, object *models.Notification) error {
	return s.publisherRepo.SendOne(ctx, object) // it calls retry inside!
}

// SendBatch sends given notifications as a batch
//
// # SendBatch is called regularly in Run
//
// might be rea-a-lly long call!
func (s *SenderService) SendBatch(ctx context.Context, objects []*models.Notification) error {
	dlqNotifications := s.publisherRepo.SendMany(ctx, objects) // it calls retry inside!

	var err error
	errGroup := &errgroup.Group{}
	errorsAmount := 0

	// if bad, then retry each one
	//
	// ... each one with retries!
	//
	// So each one is in separate goroutine
	for obj := range dlqNotifications.Items() {
		errorsAmount += 1

		errGroup.Go(func() error {

			return func(obj *dlq.DLQItem[*models.Notification]) error {

				zlog.Logger.Error().
					Err(obj.Error()).
					Stringer("id", obj.Value().ID).
					Msg("failed to send object, trying to resend...")

				err = s.QuickSend(ctx, obj.Value())
				if err != nil {
					zlog.Logger.Error().
						Err(obj.Error()).
						Stringer("id", obj.Value().ID).
						Any("object", obj.Value()).
						Msg("failed to send object!")

					return err
				} else {
					zlog.Logger.Info().
						Stringer("id", obj.Value().ID).
						Msg("successfully sent object on second try!")
				}

				return nil
			}(obj)
		})
	}

	err = errGroup.Wait()
	if err != nil {
		return fmt.Errorf("failed to send '%d' objects, example err: %w", errorsAmount, err)
	}

	return err
}

// WhenNextFetch tells when will next fetch be performed
//
// If a new notification is created and next fetch will be after its publication_at, call QuickSend on it!
func (s *SenderService) WhenNextFetch() time.Time {
	return s.nextFetchIsAt // I hope there are no races
}

// QuickSendIfNeeded is an example SignalFunc (check what it is)
//
// You can use it!
func (s *SenderService) QuickSendIfNeeded(ctx context.Context, object *models.Notification) error {
	if object.PublicationAt.Value().Before(s.WhenNextFetch()) {
		return s.QuickSend(ctx, object)
	}
	return nil
}
