package main

import (
	"consumer_worker/internal/config"
	"context"
	"fmt"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/zlog"
	"log"
	"os"
	"os/signal"
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

	rabbitConsumer := rabbitmq.NewConsumer(
		&rabbitmq.Channel{},
		&rabbitmq.ConsumerConfig{
			Queue:     cfg.RabbitMQConfig.Queue,
			Consumer:  cfg.RabbitMQConfig.Consumer,
			AutoAck:   cfg.RabbitMQConfig.AutoAck,
			Exclusive: false,
			NoWait:    cfg.RabbitMQConfig.NoWait,
			Args:      nil,
		})

	fmt.Println(rabbitConsumer)

	// TODO: scalable consumer worker that reads and reads
}
