package main

import (
	"context"
	"fmt"
	"github.com/chempik1234/wb-l3-1/consumer_worker/internal/config"
	"github.com/chempik1234/wb-l3-1/consumer_worker/internal/connect"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx := context.Background()

	// use OS signals for graceful shutdown
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	// load config from env
	cfg, err := config.NewAppConfig("", "")
	if err != nil {
		log.Fatal(fmt.Errorf("error loading config: %w", err))
	}

	// init zlog.Logger with given LogLevel
	zlog.InitConsole()
	err = zlog.SetLevel(cfg.LogConfig.LogLevel)
	if err != nil {
		log.Fatal(fmt.Errorf("error setting log level to '%s': %w", cfg.LogConfig.LogLevel, err))
	}

	//region rabbitMQ
	var rabbitConsumerConsole *rabbitmq.Consumer
	var rabbitConsumerTelegram *rabbitmq.Consumer
	var rabbitConsumerEmail *rabbitmq.Consumer
	var rabbitmqChannelToClose *rabbitmq.Channel
	rabbitConsumerConsole, rabbitmqChannelToClose, err = connect.GetRabbitMQConsumer(
		&rabbitmq.ConsumerConfig{
			Queue:     cfg.RabbitMQConfig.QueueForChannel.Console,
			Consumer:  cfg.RabbitMQConfig.Consumer,
			AutoAck:   cfg.RabbitMQConfig.AutoAck,
			Exclusive: false,
			NoWait:    cfg.RabbitMQConfig.NoWait,
			Args:      nil,
		},
		retry.Strategy{
			Attempts: cfg.RabbitMQRetryConfig.Attempts,
			Delay:    time.Duration(cfg.RabbitMQRetryConfig.DelayMilliseconds) * time.Millisecond,
			Backoff:  cfg.RabbitMQRetryConfig.Backoff,
		},
	)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("error creating rabbitmq consumer")
	}

	defer func(rabbitmqChannelToClose *rabbitmq.Channel) {
		closeErr := rabbitmqChannelToClose.Close()
		if closeErr != nil {
			zlog.Logger.Error().Err(closeErr).Msg("error closing rabbitmq channel")
		}
	}(rabbitmqChannelToClose)
	//endregion

	fmt.Println(rabbitConsumerConsole, rabbitConsumerTelegram, rabbitConsumerEmail)

	// TODO: scalable consumer worker that reads and reads
}
