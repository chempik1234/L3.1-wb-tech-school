package repositories

import (
	"context"
	"delayed_notifier/internal/dto"
	"delayed_notifier/internal/models"
	"delayed_notifier/pkg/dlq"
	"fmt"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
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
			}

		}
		DLQ.Close()
	}()

	return DLQ
}

func (n *NotificationRabbitMQ) routingKey(notification *models.Notification) string {
	return notification.Channel.String()
}
