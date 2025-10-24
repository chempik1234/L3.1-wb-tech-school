package senders

import (
	"context"
	"fmt"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/models"
	"github.com/wb-go/wbf/zlog"
)

// ConsoleSender is a sender that simply logs out the message
type ConsoleSender struct{}

// Send of ConsoleSender simply logs out a message
func (s *ConsoleSender) Send(ctx context.Context, notification *models.Notification) error {
	zlog.Logger.Log().Msg(fmt.Sprintf("%s: %s", notification.Content.Title, notification.Content.Message))
	return nil
}

// NewConsoleSender creates a new ConsoleSender
//
// requires nothing!
func NewConsoleSender() *ConsoleSender {
	return &ConsoleSender{}
}
