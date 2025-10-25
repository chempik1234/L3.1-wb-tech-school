package repositories

import (
	"context"
	"fmt"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/internal/dto"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/internal/models"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/pkg/dlq"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

// NotificationRabbitMQ is the RabbitMQ implementation on ports.NotificationPublisherRepository
type NotificationRabbitMQ struct {
	publisher     *rabbitmq.Publisher
	retryStrategy retry.Strategy
}

// NewNotificationRabbitMQ creates a new NotificationRabbitMQ
func NewNotificationRabbitMQ(
	publisher *rabbitmq.Publisher,
	retryStrategy retry.Strategy,
) *NotificationRabbitMQ {
	return &NotificationRabbitMQ{
		publisher,
		retryStrategy,
	}
}

// SendOne sends 1 notification at a time
func (n *NotificationRabbitMQ) SendOne(ctx context.Context, notification *models.Notification) error {
	body, err := dto.NotificationSendBodyFromEntityBytes(notification)
	if err != nil {
		return fmt.Errorf("couldn't create body to send one: %w", err)
	}
	err = n.publisher.PublishWithRetry(body, n.routingKey(notification), "application/json", n.retryStrategy)
	if err != nil {
		return fmt.Errorf("couldn't send message to rabbitMQ: %w", err)
	}
	zlog.Logger.Debug().Msg("sent one message to rabbitMQ")
	return nil
}

// SendMany sends batch of notifications at a time, create notifications lists and use it
func (n *NotificationRabbitMQ) SendMany(ctx context.Context, notifications []*models.Notification) *dlq.DLQ[*models.Notification] {
	DLQ := dlq.NewDLQ[*models.Notification](len(notifications) / 10)

	go func() {
		for _, notification := range notifications {
			body, err := dto.NotificationSendBodyFromEntityBytes(notification)
			if err != nil {
				DLQ.Put(notification, fmt.Errorf("couldn't send message to rabbitMQ: %w", err))
			}

			err = n.publisher.PublishWithRetry(body, n.routingKey(notification), "application/json", n.retryStrategy)
			if err != nil {
				DLQ.Put(notification, fmt.Errorf("couldn't send message to rabbitMQ: %w", err))
			} else {
				zlog.Logger.Debug().Msg("sent message in batch to rabbitMQ")
			}
		}
		DLQ.Close()
	}()

	return DLQ
}

func (n *NotificationRabbitMQ) routingKey(notification *models.Notification) string {
	return notification.Channel.String()
}
