package connect

import (
	"delayed_notifier/internal/config"
	"fmt"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
)

// GetRabbitMQConsumer simplifies complex rabbitMQ connection process!
//
// returns:
//
//	consumer
//	channel to close
//	error
func GetRabbitMQConsumer(rabbitCfg config.RabbitMQConfig, rabbitmqRetryStrategy retry.Strategy) (*rabbitmq.Consumer, *rabbitmq.Channel, error) {
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

	queues := []string{rabbitCfg.QueueSend.Email, rabbitCfg.QueueSend.Telegram, rabbitCfg.QueueSend.Console}
	for _, queue := range queues {
		err = retry.Do(
			func() error {
				_, errQueue := rabbitMQQueueManager.DeclareQueue(queue)
				return errQueue
			},
			rabbitmqRetryStrategy,
		)

		if err != nil {
			return nil, nil, fmt.Errorf("error declaring queue '%s': %w", queue, err)
		}
	}

	// final step. create publisher
	rabbitmqPublisher := rabbitmq.NewPublisher(
		rabbitMQChannel,
		rabbitCfg.Exchange)
	return rabbitmqPublisher, rabbitMQChannel, nil
}
