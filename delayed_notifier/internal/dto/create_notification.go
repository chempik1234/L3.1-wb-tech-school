package dto

import (
	"fmt"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/internal/internaltypes"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/internal/models"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/pkg/types"
)

// CreateNotificationBody is a DTO for create endpoint
type CreateNotificationBody struct {
	PublicationAt string                  `json:"publication_at"`
	Channel       string                  `json:"channel"`
	Content       notificationBodyContent `json:"content"`
}

// ToEntity is a method that converts DTO into create-able model (without ID)
func (b CreateNotificationBody) ToEntity() (*models.Notification, error) {
	var err error

	// publication_at
	var publicationAt types.DateTime
	publicationAt, err = types.NewDateTimeFromString(b.PublicationAt)
	if err != nil {
		return nil, fmt.Errorf("incorrect 'publication_at' '%s': %w", b.PublicationAt, err)
	}

	// channel
	var channel internaltypes.NotificationChannel
	channel, err = internaltypes.NotificationChannelFromString(b.Channel)
	if err != nil {
		return nil, fmt.Errorf("incorrect 'channel' '%s': %w", b.Channel, err)
	}

	// content body
	title := types.NewAnyText(b.Content.Title)
	message := types.NewAnyText(b.Content.Message)

	// result
	return &models.Notification{
		PublicationAt: publicationAt,
		Channel:       channel,
		Content: models.NotificationContent{
			Title:   title,
			Message: message,
		},
	}, nil
}
