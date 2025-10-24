package models

import (
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/internaltypes"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/pkg/types"
)

// Notification is the main model - it's saved in DB, cached and DTOs are converted to it
type Notification struct {
	PublicationAt types.DateTime
	ID            *types.UUID
	Channel       internaltypes.NotificationChannel
	Content       NotificationContent
	Sent          bool
}

// NotificationContent is the universal struct for content: notification has a title and a message
type NotificationContent struct {
	Title   types.AnyText
	Message types.AnyText
}
