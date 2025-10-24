package dto

import (
	"encoding/json"
	"fmt"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/models"
)

// NotificationSendBody is the DTO for sending to MQ
type NotificationSendBody struct {
	Content       notificationBodyContent `json:"content"`
	ID            string                  `json:"id"`
	PublicationAt string                  `json:"publication_at"`
	Channel       string                  `json:"channel"`
}

type notificationBodyContent struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// NotificationSendBodyFromEntity creates a new *NotificationSendBody from given object
//
// Use it to send to MQ
func NotificationSendBodyFromEntity(object *models.Notification) *NotificationSendBody {
	return &NotificationSendBody{
		Content: notificationBodyContent{
			Title:   object.Content.Title.String(),
			Message: object.Content.Message.String(),
		},
		ID:            object.ID.String(),
		PublicationAt: object.PublicationAt.String(),
		Channel:       object.Channel.String(),
	}
}

// NotificationSendBodyFromEntityBytes creates a ready-to-send []byte body from given object
//
// uses NotificationSendBodyFromEntity
func NotificationSendBodyFromEntityBytes(object *models.Notification) ([]byte, error) {
	result, err := json.Marshal(NotificationSendBodyFromEntity(object))
	if err != nil {
		return nil, fmt.Errorf("could not marshal NotificationSendBody: %w", err)
	}
	return result, nil
}
