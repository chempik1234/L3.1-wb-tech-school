package connect

import (
	"fmt"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/internaltypes"
	"github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
)

// RabbitMQConsumerConfig is the options struct for GetRabbitMQConsumer
//
// too much fields to pass as arguments in function
type RabbitMQConsumerConfig struct {
	Exchange string
	User     string
	Password string
	Host     string
	Port     int
	VHost    string

	QueueName string
	Consumer  string
	AutoAck   bool
	NoWait    bool
}

// GetRabbitMQConsumer simplifies complex rabbitMQ connection process!
//
// returns:
//
//	consumer
//	channel to close
//	error
func GetRabbitMQConsumer(rabbitCfg RabbitMQConsumerConfig, rabbitmqRetryStrategy retry.Strategy) (*rabbitmq.Consumer, *rabbitmq.Channel, error) {
	// step 1. init connect
	rabbitMQConn, err := rabbitmq.Connect(
		fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
			rabbitCfg.User,
			rabbitCfg.Password,
			rabbitCfg.Host,
			rabbitCfg.Port,
			rabbitCfg.VHost,
		), rabbitmqRetryStrategy.Attempts, rabbitmqRetryStrategy.Delay)
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to rabbitmq: %w", err)
	}

	// step 2. get channel to bind
	var rabbitMQChannel *rabbitmq.Channel
	rabbitMQChannel, err = rabbitMQConn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to rabbitmq (conn.Channel()): %w", err)
	}

	// step 3. bind channel to exchange with type direct
	rabbitMQExchange := rabbitmq.NewExchange(rabbitCfg.Exchange, "direct")
	err = rabbitMQExchange.BindToChannel(rabbitMQChannel)
	if err != nil {
		return nil, nil, fmt.Errorf("error binding rabbitmq channel to exchange '%s': %w",
			rabbitCfg.Exchange, err)
	}

	// step 4. declare queues (at least try)
	rabbitMQQueueManager := rabbitmq.NewQueueManager(rabbitMQChannel)

	// bind to the queue by using channel names as routing keys!
	err = retry.Do(
		func() error {
			q, errQueue := rabbitMQQueueManager.DeclareQueue(rabbitCfg.QueueName)
			if errQueue == nil {
				for _, channelString := range internaltypes.ChannelAllStrings {
					errQueue = rabbitMQChannel.QueueBind(q.Name, channelString, rabbitMQExchange.Name(), rabbitCfg.NoWait, make(amqp091.Table))
					if errQueue != nil {
						break
					}
				}
			}
			return errQueue
		},
		rabbitmqRetryStrategy,
	)

	if err != nil {
		return nil, nil, fmt.Errorf("error declaring queue '%s': %w", rabbitCfg.QueueName, err)
	}

	// final step. create consumer
	rabbitmqPublisher := rabbitmq.NewConsumer(
		rabbitMQChannel,
		rabbitmq.NewConsumerConfig(rabbitCfg.QueueName))
	return rabbitmqPublisher, rabbitMQChannel, nil
}
