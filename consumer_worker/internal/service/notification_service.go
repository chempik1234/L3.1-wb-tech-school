package service

import (
	"container/heap"
	"context"
	"errors"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/internaltypes"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/models"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/notificationheap"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/ports"
	"github.com/wb-go/wbf/zlog"
	"sync"
	"time"
)

var ErrUnknownChannel = errors.New("unknown channel: no sender for it")

// NotificationService is the main service that reads, sorts and sends 100 MLN notifications per 1 MS
//
//	s.StartReceiving(ctx)
//	<-ctx.Done()
//	// service already stopped
type NotificationService struct {
	// no writes -> no mutex
	channelToSender map[internaltypes.NotificationChannel]ports.NotificationSender

	receiver ports.NotificationReceiver

	notificationHeap *notificationheap.NotificationHeap
	heapMutex        sync.RWMutex

	checkPeriod time.Duration
}

// NewNotificationService creates a new NotificationService
func NewNotificationService(receiver ports.NotificationReceiver, channelToSender map[internaltypes.NotificationChannel]ports.NotificationSender, checkPeriod time.Duration) *NotificationService {
	// make container/heap handle our sorting
	notificationHeap := &notificationheap.NotificationHeap{}
	heap.Init(notificationHeap)

	return &NotificationService{
		channelToSender:  channelToSender,
		receiver:         receiver,
		heapMutex:        sync.RWMutex{},
		notificationHeap: notificationHeap,
		checkPeriod:      checkPeriod,
	}
}

// Run is the life cycle function
//
// blocking, stops automatically with ctx.Done()
func (s *NotificationService) Run(ctx context.Context) error {
	objects := s.receiver.StartReceiving()
	var object *models.Notification

	go s.serveHeap(ctx)

out:
	for {
		select {
		case <-ctx.Done():
			break out
		case object = <-objects:
			break
		}

		// check if we'll be able to even send this notification
		if _, ok := s.channelToSender[object.Channel]; !ok {
			zlog.Logger.Error().Stringer("channel", &object.Channel).Msg("unable to find sender for given channel")
			continue
		}

		s.heapMutex.Lock()
		heap.Push(s.notificationHeap, object)
		s.heapMutex.Unlock()
	}

	return s.receiver.StopReceiving()
}

// serveHeap simply reads everything from heap every 50ms and sends it
func (s *NotificationService) serveHeap(ctx context.Context) {
	// step 1. Peek
	// step 2. Pop
	// step 3. Send

	ticker := time.NewTicker(s.checkPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()

			s.heapMutex.Lock()

			// read everything from heap
			for s.notificationHeap.Len() > 0 {
				// step 1. Peek notification for some validation
				notificationToPublish := s.notificationHeap.Peek()
				if notificationToPublish == nil {
					break
				}

				// step 1.1. If time isn't even soon, we leave
				if publicationTime := notificationToPublish.PublicationAt.Value(); publicationTime.Add(s.checkPeriod).After(now) {
					break
				}

				// step 2. Pop
				notification := heap.Pop(s.notificationHeap).(*models.Notification)

				s.heapMutex.Unlock()

				// step 3. Send
				if err := s.sendNotification(ctx, notification); err != nil {
					zlog.Logger.Error().
						Err(err).
						Str("notification_id", notification.ID.String()).
						Str("channel", notification.Channel.String()).
						Msg("failed to send notification")
				} else {
					zlog.Logger.Info().
						Str("notification_id", notification.ID.String()).
						Str("channel", notification.Channel.String()).
						Msg("notification sent successfully")
				}

				s.heapMutex.Lock()
			}

			s.heapMutex.Unlock()
		}
	}
}

func (s *NotificationService) sendNotification(ctx context.Context, notification *models.Notification) error {
	sender, ok := s.channelToSender[notification.Channel]
	if !ok {
		return ErrUnknownChannel
	}
	err := sender.Send(ctx, notification)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("notification_id", notification.ID.String()).
			Msg("failed to send notification")
	}
	return err
}
