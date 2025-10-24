package ports

import (
	"context"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/models"
)

// NotificationReceiver is the port for Notification receiver
//
// Used in the service as messages source
type NotificationReceiver interface {
	// Receive reads 1 message and returns it
	// Receive(ctx context.Context) (*models.Notification, error)

	// StartReceiving begins the consuming and returns readonly channel with parsed objects
	StartReceiving() <-chan *models.Notification

	// StopReceiving stops the consuming, must be called in the end
	StopReceiving() error
}

// NotificationSender is the port for sender
//
// Used in the service as messages sender variation
type NotificationSender interface {
	// Send sends a message to whatever the Implementation is created for
	Send(ctx context.Context, notification *models.Notification) error
}
