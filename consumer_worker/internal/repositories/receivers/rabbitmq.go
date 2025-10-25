package receivers

import (
	"encoding/json"
	"fmt"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/dto"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/internaltypes"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/models"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/pkg/types"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

// RabbitMQReceiver is the ports.NotificationReceiver Repository for RabbitMQ
//
//	defer channel.Close() // or whatever
//	defer r.StopReceiving()
//
//	err = r.StartReceiving()
//	// handle err
//
//	//...
//
//	obj, err = r.Receive(ctx)
type RabbitMQReceiver struct {
	consumer      *rabbitmq.Consumer
	channel       *rabbitmq.Channel
	retryStrategy retry.Strategy

	// close messages - objects close too
	// StartReceiving launches transfer from 1 to 2 channel

	messages    chan []byte
	objectsChan chan *models.Notification
}

// NewRabbitMQReceiver creates a new RabbitMQReceiver for given consumer
func NewRabbitMQReceiver(consumer *rabbitmq.Consumer, channel *rabbitmq.Channel, retryStrategy retry.Strategy) *RabbitMQReceiver {
	return &RabbitMQReceiver{
		consumer:      consumer,
		channel:       channel,
		messages:      make(chan []byte),
		objectsChan:   make(chan *models.Notification),
		retryStrategy: retryStrategy,
	}
}

// StartReceiving starts the consuming, in background
//
// Must be called
func (r *RabbitMQReceiver) StartReceiving() <-chan *models.Notification {
	go func() {
		err := r.consumer.ConsumeWithRetry(r.messages, r.retryStrategy)
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("error occurred while consuming messages")
		}
	}()

	go func() {
		defer close(r.objectsChan)
		for delivery := range r.messages {
			object, err := r.processMessage(delivery)
			if err != nil {
				zlog.Logger.Info().Err(err).Msg("error while processing message")
				continue
			}

			r.objectsChan <- object
		}
	}()

	return r.objectsChan
}

// StopReceiving stops the processing of messages.
//
// Must be called
func (r *RabbitMQReceiver) StopReceiving() error {
	err := r.channel.Close()
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("error closing rabbitmq channel")
	}
	close(r.messages)

	return err
}

// processMessage reads 1 message from RabbitMQ and converts it to Notification model
func (r *RabbitMQReceiver) processMessage(delivery []byte) (*models.Notification, error) {
	// step 1. read (passed)

	// step 2. parse
	var messageData dto.NotificationSendBody

	// step 2.1: no ack if bad content! but wbf limits me
	if err := json.Unmarshal(delivery, &messageData); err != nil {
		return nil, fmt.Errorf("bad message (bad json): %w", err)
	}

	// step 3. deserialize
	notification, err := dto.NotificationModelFromSendDTO(&messageData)
	if err != nil {
		return nil, fmt.Errorf("bad message (could't convert to model): %w", err)
	}

	zlog.Logger.Debug().
		Str("notification_id", notification.ID.String()).
		Str("channel", notification.Channel.String()).
		Msg("notification received from rabbitmq")

	return notification, nil
}

func (r *RabbitMQReceiver) convertToNotification(data struct {
	Content struct {
		Title   string `json:"title"`
		Message string `json:"message"`
	} `json:"content"`
	ID            string `json:"id"`
	PublicationAt string `json:"publication_at"`
	Channel       string `json:"channel"`
}) (*models.Notification, error) {
	uuid, err := types.NewUUID(data.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %w", err)
	}

	channel, err := internaltypes.NotificationChannelFromString(data.Channel)
	if err != nil {
		return nil, fmt.Errorf("invalid channel: %w", err)
	}

	publicationAt, err := types.NewDateTimeFromString(data.PublicationAt)
	if err != nil {
		return nil, fmt.Errorf("invalid publication date: %w", err)
	}

	content := models.NotificationContent{
		Title:   types.NewAnyText(data.Content.Title),
		Message: types.NewAnyText(data.Content.Message),
	}

	return &models.Notification{
		ID:            &uuid,
		Channel:       channel,
		PublicationAt: publicationAt,
		Content:       content,
		Sent:          false,
	}, nil
}
