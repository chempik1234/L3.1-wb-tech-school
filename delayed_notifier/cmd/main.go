package main

import (
	"context"
	"delayed_notifier/internal/config"
	"delayed_notifier/internal/connect"
	"delayed_notifier/internal/repositories"
	"delayed_notifier/internal/service"
	"delayed_notifier/internal/transport"
	"delayed_notifier/pkg/http_server"
	"fmt"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
	"log"
	"sync"
	"time"
)

func main() {
	//region load config from env
	cfg, err := config.NewAppConfig("/app/config.yaml", "")
	if err != nil {
		log.Fatal(fmt.Errorf("error loading config: %w", err))
	}
	//endregion

	//region init zlog.Logger with given LogLevel
	zlog.InitConsole()
	err = zlog.SetLevel(cfg.LogConfig.LogLevel)
	if err != nil {
		zlog.Logger.Fatal().Err(fmt.Errorf("error setting log level to '%s': %w", cfg.LogConfig.LogLevel, err))
	}
	//endregion

	//region retry (define first for later postgres, rabbitmq, redis connections)
	postgresRetryStategy := retry.Strategy{
		Attempts: cfg.PostgresRetryConfig.Attempts,
		Delay:    time.Duration(cfg.PostgresRetryConfig.DelayMilliseconds) * time.Millisecond,
		Backoff:  cfg.PostgresRetryConfig.Backoff,
	}

	redisRetryStategy := retry.Strategy{
		Attempts: cfg.RedisRetryConfig.Attempts,
		Delay:    time.Duration(cfg.RedisRetryConfig.DelayMilliseconds) * time.Millisecond,
		Backoff:  cfg.RedisRetryConfig.Backoff,
	}

	rabbitmqRetryStrategy := retry.Strategy{
		Attempts: cfg.RabbitMQRetryConfig.Attempts,
		Delay:    time.Duration(cfg.RabbitMQRetryConfig.DelayMilliseconds) * time.Millisecond,
		Backoff:  cfg.RabbitMQRetryConfig.Backoff,
	}

	zlog.Logger.Info().Msg("retry policies created")
	//endregion

	//region rabbitMQ
	var rabbitmqPublisher *rabbitmq.Publisher
	var rabbitmqChannelToClose *rabbitmq.Channel
	rabbitmqPublisher, rabbitmqChannelToClose, err = connect.GetRabbitMQPublisher(cfg.RabbitMQConfig, rabbitmqRetryStrategy)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("error creating rabbitmq publisher")
	}

	defer func(rabbitmqChannelToClose *rabbitmq.Channel) {
		closeErr := rabbitmqChannelToClose.Close()
		if closeErr != nil {
			zlog.Logger.Error().Err(closeErr).Msg("error closing rabbitmq channel")
		}
	}(rabbitmqChannelToClose)

	zlog.Logger.Info().Msg("rabbitMQ publisher created")
	//endregion

	//region postgres
	var postgresDB *dbpg.DB

	// connect to postgres with retry
	err = retry.Do(
		func() error {
			var postgresConnErr error

			postgresDB, postgresConnErr = dbpg.New(

				cfg.PostgresConfig.MasterDSN,
				cfg.PostgresConfig.SlaveDSNs,

				&dbpg.Options{
					MaxOpenConns:    cfg.PostgresConfig.MaxOpenConnections,
					MaxIdleConns:    cfg.PostgresConfig.MaxIdleConnections,
					ConnMaxLifetime: time.Duration(cfg.PostgresConfig.ConnectionMaxLifetimeSeconds) * time.Second,
				})

			return postgresConnErr
		},
		postgresRetryStategy)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("couldn't create postgres balancer")
	}

	zlog.Logger.Info().Msg("postgres balancer created")
	//endregion

	//region redis
	redisClient := redis.New(
		cfg.RedisConfig.Addr,
		cfg.RedisConfig.Password,
		cfg.RedisConfig.DB,
	)
	redisExpiration := time.Second * time.Duration(cfg.RedisConfig.TTLSeconds)
	zlog.Logger.Info().Msg("redis created")
	//endregion

	//region services
	postgresRepo := repositories.NewNotificationPostgres(postgresDB, postgresRetryStategy)
	redisRepo := repositories.NewNotificationRedis(redisClient, redisRetryStategy, redisExpiration)
	rabbitmqRepo := repositories.NewNotificationRabbitMQ(rabbitmqPublisher, rabbitmqRetryStrategy)

	senderService := service.NewSenderService(
		time.Duration(cfg.FetcherConfig.FetchPeriodSeconds)*time.Second,
		time.Duration(cfg.FetcherConfig.FetchMaxDiapasonSeconds)*time.Second,
		rabbitmqRepo, postgresRepo,
	)
	crudService := service.NewNotificationCRUDService(postgresRepo, redisRepo, senderService.QuickSendIfNeeded)
	//endregion

	ctx, stopCtx := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	//region run in background
	wg.Add(1)
	go func(wg *sync.WaitGroup, ctx2 context.Context) {
		defer wg.Done()
		senderService.Run(ctx)
	}(wg, ctx)
	//endregion

	//region Start HTTP
	notifyHTTPHandler := transport.NewNotifyHandler(crudService)
	appRouter := transport.AssembleRouter(notifyHTTPHandler)
	appServer := http_server.NewHTTPServer(appRouter)

	zlog.Logger.Info().Int("http_port", cfg.ServerConfig.HTTPPort).Msg("server starting :http_port")

	err = appServer.GracefulRun(ctx, cfg.ServerConfig.HTTPPort)
	//endregion

	//region shutdown
	if err != nil {
		zlog.Logger.Error().Msg(fmt.Errorf("http server error: %w", err).Error())
	}

	zlog.Logger.Info().Msg("server gracefully stopped")

	stopCtx()
	wg.Wait()
	zlog.Logger.Info().Msg("background operations gracefully stopped")
	//endregion
}
