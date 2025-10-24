package connect

import (
	"fmt"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/internal/config"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
)

// GetRabbitMQPublisher simplifies complex rabbitMQ connection process!
//
// returns:
//
//	publisher
//	channel to close
//	error
func GetRabbitMQPublisher(rabbitCfg config.RabbitMQConfig, rabbitmqRetryStrategy retry.Strategy) (*rabbitmq.Publisher, *rabbitmq.Channel, error) {
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

	err = retry.Do(
		func() error {
			_, errQueue := rabbitMQQueueManager.DeclareQueue(rabbitCfg.QueueSend)
			return errQueue
		},
		rabbitmqRetryStrategy,
	)

	if err != nil {
		return nil, nil, fmt.Errorf("error declaring queue '%s': %w", rabbitCfg.QueueSend, err)
	}

	// final step. create publisher
	rabbitmqPublisher := rabbitmq.NewPublisher(
		rabbitMQChannel,
		rabbitCfg.Exchange)
	return rabbitmqPublisher, rabbitMQChannel, nil
}
