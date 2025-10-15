package dto

import (
	"delayed_notifier/internal/models"
)

// FullNotificationBody is a DTO for fully-serialized Notification model
type FullNotificationBody struct {
	Content       notificationBodyContent `json:"content"`
	ID            string                  `json:"id"`
	PublicationAt string                  `json:"publication_at"`
	Channel       string                  `json:"channel"`
	Sent          bool                    `json:"sent"`
}

type notificationBodyContent struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// FullNotificationBodyFromEntity is a method that model to DTO, used for “return “
func FullNotificationBodyFromEntity(model *models.Notification) *FullNotificationBody {
	return &FullNotificationBody{
		ID:            model.ID.String(),
		PublicationAt: model.PublicationAt.String(),
		Channel:       model.Channel.String(),
		Content: notificationBodyContent{
			Title:   model.Content.Title.String(),
			Message: model.Content.Message.String(),
		},
		Sent: model.Sent,
	}
}
