package senders

import (
	"context"
	"fmt"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/models"
)

// ConsoleSender is a sender that simply logs out the message
type ConsoleSender struct{}

// Send of ConsoleSender simply logs out a message with several fmt.Printf
func (s *ConsoleSender) Send(ctx context.Context, notification *models.Notification) error {
	// zlog.Logger.Log().Msg(fmt.Sprintf("%s: %s", notification.Content.Title, notification.Content.Message))
	fmt.Printf("=== SEND CONSOLE NOTIFICATION ===\n")
	fmt.Printf("ID: %s\n", notification.ID.String())
	fmt.Printf("Channel: %s\n", notification.Channel.String())
	fmt.Printf("Title: %s\n", notification.Content.Title.String())
	fmt.Printf("Message: %s\n", notification.Content.Message.String())
	fmt.Printf("Publication At: %s\n", notification.PublicationAt.String())
	fmt.Printf("=================================\n\n")
	return nil
}

// NewConsoleSender creates a new ConsoleSender
//
// requires nothing!
func NewConsoleSender() *ConsoleSender {
	return &ConsoleSender{}
}
