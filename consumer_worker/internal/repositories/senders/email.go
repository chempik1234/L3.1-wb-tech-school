package senders

import (
	"context"
	"fmt"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/models"
	"github.com/wb-go/wbf/retry"
	"net/smtp"
)

// EmailSender is a sender that sends an email with retries
type EmailSender struct {
	from     string
	password string
	// host:port
	addr string
	host string

	retryStrategy retry.Strategy
}

// Send of EmailSender simply logs out a message with net/smtp
func (s *EmailSender) Send(ctx context.Context, notification *models.Notification) error {
	err := retry.Do(
		func() error { return s.sendMail(ctx, notification) },
		s.retryStrategy)
	if err != nil {
		return fmt.Errorf("error send email: %w", err)
	}
	return nil
}

func (s *EmailSender) sendMail(ctx context.Context, notification *models.Notification) error {
	to := []string{notification.SendTo.String()}
	from := s.from

	msg := fmt.Sprintf("%s\n\n%s", notification.Content.Title, notification.Content.Message)

	return smtp.SendMail(
		s.addr,
		smtp.PlainAuth("", from, s.password, s.host),
		from, to,
		[]byte(msg),
	)
}

// NewEmailSender creates a new EmailSender
func NewEmailSender(fromMail, password, host string, port int, retryStrategy retry.Strategy) *EmailSender {
	return &EmailSender{
		from:          fromMail,
		password:      password,
		addr:          fmt.Sprintf("%s:%d", host, port),
		retryStrategy: retryStrategy,
		host:          host,
	}
}
