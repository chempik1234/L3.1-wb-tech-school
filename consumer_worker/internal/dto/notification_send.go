package dto

import (
	"fmt"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/internaltypes"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/models"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/pkg/types"
)

// NotificationSendBody is the DTO for sending to MQ
type NotificationSendBody struct {
	Content       notificationBodyContent `json:"content"`
	ID            string                  `json:"id"`
	PublicationAt string                  `json:"publication_at"`
	Channel       string                  `json:"channel"`
	SendTo        string                  `json:"send_to,omitempty"`
}

type notificationBodyContent struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// NotificationModelFromSendDTO deserializes DTO into a normal *models.Notification struct
func NotificationModelFromSendDTO(dto *NotificationSendBody) (*models.Notification, error) {
	publicationAt, err := types.NewDateTimeFromString(dto.PublicationAt)
	if err != nil {
		return nil, fmt.Errorf("invalid publication_at: %w", err)
	}

	id, err := types.NewUUID(dto.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	channel, err := internaltypes.NotificationChannelFromString(dto.Channel)
	if err != nil {
		return nil, fmt.Errorf("invalid channel: %w", err)
	}

	sendTo, err := internaltypes.NewSendTo(types.NewAnyText(dto.SendTo), channel)
	if err != nil {
		return nil, fmt.Errorf("invalid send_to: %w", err)
	}

	return &models.Notification{
		PublicationAt: publicationAt,
		ID:            &id,
		Channel:       channel,
		Content: models.NotificationContent{
			Title:   types.NewAnyText(dto.Content.Title),
			Message: types.NewAnyText(dto.Content.Message),
		},
		Sent:   true,
		SendTo: sendTo,
	}, nil
}
