package ports

import (
	"context"
	"delayed_notifier/internal/models"
	"delayed_notifier/pkg/dlq"
	"delayed_notifier/pkg/types"
)

// NotificationFetcherRepository is the port for fetching a batch-to-send from DB & changing status
type NotificationFetcherRepository interface {
	// Fetch fetches objects to be sent (only up to maxPublicationAt not to store everything in memory)=
	//
	// fetches are supposed to be done regularly
	Fetch(ctx context.Context, maxPublicationAt types.DateTime) ([]*models.Notification, error)

	// MarkAsSent should be used to mark fetched notifications as sent
	MarkAsSent(ctx context.Context, ids []*types.UUID) error
}

// NotificationPublisherRepository is the port for notification sender
//
// They are supposed to be sent with a MQ, such as RabbitMQ, to distributed sender workers
type NotificationPublisherRepository interface {
	// SendOne sends 1 notification at a time
	SendOne(ctx context.Context, notification *models.Notification) error

	// SendMany sends batch of notifications at a time, create notifications lists and use it
	SendMany(ctx context.Context, notifications []*models.Notification) *dlq.DLQ[*models.Notification]
}
