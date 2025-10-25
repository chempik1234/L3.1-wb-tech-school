package config

import (
	"fmt"
	"github.com/wb-go/wbf/config"
)

// AppConfig is THE whole config struct
type AppConfig struct {
	LogConfig           LogConfig           `env-prefix:"LOG_"`
	RabbitMQConfig      RabbitMQConfig      `env-prefix:"RABBITMQ_"`
	RabbitMQRetryConfig RetryStrategyConfig `env-prefix:"RETRY_RABBITMQ_"`
}

// NewAppConfig creates a new struct of "THE config"
func NewAppConfig(configFilePath, envFilePath string) (*AppConfig, error) {
	appConfig := &AppConfig{}

	cfg := config.New()

	//region defaults
	cfg.SetDefault("consumer_worker.log.level", "info")

	cfg.SetDefault("consumer_worker.retry_rabbitmq.attempts", 3)
	cfg.SetDefault("consumer_worker.retry_rabbitmq.delay_milliseconds", 300)
	cfg.SetDefault("consumer_worker.retry_rabbitmq.backoff", 1.5)
	//endregion

	// region flags

	// why flags lol dude
	cfg.ParseFlags()
	//endregion

	err := cfg.Load(configFilePath, envFilePath, "")
	if err != nil {
		return appConfig, fmt.Errorf("failed to load config: %w", err)
	}

	// LogConfig
	appConfig.LogConfig.LogLevel = cfg.GetString("consumer_worker.log.level")

	// RabbitMQConfig
	appConfig.RabbitMQConfig.Exchange = cfg.GetString("consumer_worker.rabbitmq.exchange")
	appConfig.RabbitMQConfig.User = cfg.GetString("consumer_worker.rabbitmq.user")
	appConfig.RabbitMQConfig.Password = cfg.GetString("consumer_worker.rabbitmq.password")
	appConfig.RabbitMQConfig.Host = cfg.GetString("consumer_worker.rabbitmq.host")
	appConfig.RabbitMQConfig.Port = cfg.GetInt("consumer_worker.rabbitmq.port")
	appConfig.RabbitMQConfig.VHost = cfg.GetString("consumer_worker.rabbitmq.vhost")
	appConfig.RabbitMQConfig.UniversalQueue = cfg.GetString("consumer_worker.rabbitmq.queue")
	// appConfig.RabbitMQConfig.QueueForChannel.Telegram = cfg.GetString("consumer_worker.rabbitmq.queue_read.telegram")
	// appConfig.RabbitMQConfig.QueueForChannel.Console = cfg.GetString("consumer_worker.rabbitmq.queue_read.console")

	// Retries
	appConfig.RabbitMQRetryConfig.Attempts = cfg.GetInt("consumer_worker.retry_rabbitmq.attempts")
	appConfig.RabbitMQRetryConfig.DelayMilliseconds = cfg.GetInt("consumer_worker.retry_rabbitmq.delay_milliseconds")
	appConfig.RabbitMQRetryConfig.Backoff = cfg.GetFloat64("consumer_worker.retry_rabbitmq.backoff")

	return appConfig, nil
}
