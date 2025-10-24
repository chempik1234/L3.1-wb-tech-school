package main

import (
	"context"
	"fmt"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/config"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/connect"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/internaltypes"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/ports"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/repositories/receivers"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/repositories/senders"
	"github.com/chempik1234/L3.1-wb-tech-school/consumer_worker/internal/service"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	ctx := context.Background()

	// use OS signals for graceful shutdown
	ctx, ctxStop := signal.NotifyContext(ctx, os.Interrupt)
	defer ctxStop()

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

	// this var is going to be changed for each channel
	rabbitConnectCfg := connect.RabbitMQConsumerConfig{
		Exchange:  cfg.RabbitMQConfig.Exchange,
		User:      cfg.RabbitMQConfig.User,
		Password:  cfg.RabbitMQConfig.Password,
		Host:      cfg.RabbitMQConfig.Host,
		Port:      cfg.RabbitMQConfig.Port,
		VHost:     cfg.RabbitMQConfig.VHost,
		QueueName: cfg.RabbitMQConfig.UniversalQueue,
		Consumer:  cfg.RabbitMQConfig.Consumer,
		AutoAck:   cfg.RabbitMQConfig.AutoAck,
		NoWait:    cfg.RabbitMQConfig.NoWait,
	}

	rabbitmqRetryStrategy := retry.Strategy{
		Attempts: cfg.RabbitMQRetryConfig.Attempts,
		Delay:    time.Duration(cfg.RabbitMQRetryConfig.DelayMilliseconds) * time.Millisecond,
		Backoff:  cfg.RabbitMQRetryConfig.Backoff,
	}

	var rabbitConsumer *rabbitmq.Consumer
	var rabbitmqChannelToClose *rabbitmq.Channel
	rabbitConsumer, rabbitmqChannelToClose, err = connect.GetRabbitMQConsumer(
		rabbitConnectCfg,
		rabbitmqRetryStrategy,
	)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("error creating rabbitmq consumer")
	}
	//endregion

	//region service
	rabbitmqReceiver := receivers.NewRabbitMQReceiver(rabbitConsumer, rabbitmqChannelToClose, rabbitmqRetryStrategy)

	channelToSender := map[internaltypes.NotificationChannel]ports.NotificationSender{
		internaltypes.ChannelConsole:  senders.NewConsoleSender(),
		internaltypes.ChannelEmail:    senders.NewConsoleSender(), // TODO: change
		internaltypes.ChannelTelegram: senders.NewConsoleSender(), // TODO: change
	}

	notificationService := service.NewNotificationService(rabbitmqReceiver, channelToSender)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		errService := notificationService.Run(ctx)
		if errService != nil {
			zlog.Logger.Error().Err(errService).Msg("error running notification service")
		}
	}(wg)
	//endregion

	<-ctx.Done()

	// notificationService stops with ctx

	wg.Wait()
	zlog.Logger.Info().Msg("shutdown complete")
}
