package service

import (
	"context"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/internaltypes"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/models"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/ports"
	"github.com/wb-go/wbf/zlog"
)

// NotificationService is the main service that reads, sorts and sends 100 MLN notifications per 1 MS
//
//	s.StartReceiving(ctx)
//	<-ctx.Done()
//	// service already stopped
type NotificationService struct {
	channelToSender map[internaltypes.NotificationChannel]ports.NotificationSender
	receiver        ports.NotificationReceiver
}

// NewNotificationService creates a new NotificationService
func NewNotificationService(receiver ports.NotificationReceiver, channelToSender map[internaltypes.NotificationChannel]ports.NotificationSender) *NotificationService {
	return &NotificationService{
		channelToSender: channelToSender,
		receiver:        receiver,
	}
}

// Run is the life cycle function
//
// blocking, stops automatically with ctx.Done()
func (s *NotificationService) Run(ctx context.Context) error {
	objects := s.receiver.StartReceiving()
	var object *models.Notification
out:
	for {
		select {
		case <-ctx.Done():
			break out
		case object = <-objects:
			break
		}

		sender, ok := s.channelToSender[object.Channel]
		if !ok {
			zlog.Logger.Error().Stringer("channel", &object.Channel).Msg("unable to find sender for given channel")
			continue
		}

		// TODO: min-heap & timer
		err := sender.Send(ctx, object)
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("unable to send notification")
		}
	}

	return s.receiver.StopReceiving()
}
