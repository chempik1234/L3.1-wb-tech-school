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
	cfg.SetDefault("delayed_notifier.log.level", "info")

	cfg.SetDefault("delayed_notifier.retry_rabbitmq.attempts", 3)
	cfg.SetDefault("delayed_notifier.retry_rabbitmq.delay_milliseconds", 300)
	cfg.SetDefault("delayed_notifier.retry_rabbitmq.backoff", 1.5)
	//endregion

	// region flags

	// why flags lol dude
	cfg.ParseFlags()
	//endregion

	err := cfg.Load(configFilePath, envFilePath, "CONSUMER_WORKER_")
	if err != nil {
		return appConfig, fmt.Errorf("failed to load config: %w", err)
	}

	// LogConfig
	appConfig.LogConfig.LogLevel = cfg.GetString("delayed_notifier.log.level")

	// RabbitMQConfig
	appConfig.RabbitMQConfig.Exchange = cfg.GetString("delayed_notifier.rabbitmq.exchange")
	appConfig.RabbitMQConfig.User = cfg.GetString("delayed_notifier.rabbitmq.user")
	appConfig.RabbitMQConfig.Password = cfg.GetString("delayed_notifier.rabbitmq.password")
	appConfig.RabbitMQConfig.Host = cfg.GetString("delayed_notifier.rabbitmq.host")
	appConfig.RabbitMQConfig.Port = cfg.GetInt("delayed_notifier.rabbitmq.port")
	appConfig.RabbitMQConfig.VHost = cfg.GetString("delayed_notifier.rabbitmq.vhost")
	appConfig.RabbitMQConfig.UniversalQueue = cfg.GetString("delayed_notifier.rabbitmq.queue")
	// appConfig.RabbitMQConfig.QueueForChannel.Telegram = cfg.GetString("delayed_notifier.rabbitmq.queue_read.telegram")
	// appConfig.RabbitMQConfig.QueueForChannel.Console = cfg.GetString("delayed_notifier.rabbitmq.queue_read.console")

	// Retries
	appConfig.RabbitMQRetryConfig.Attempts = cfg.GetInt("delayed_notifier.retry_rabbitmq.attempts")
	appConfig.RabbitMQRetryConfig.DelayMilliseconds = cfg.GetInt("delayed_notifier.retry_rabbitmq.delay_milliseconds")
	appConfig.RabbitMQRetryConfig.Backoff = cfg.GetFloat64("delayed_notifier.retry_rabbitmq.backoff")

	return appConfig, nil
}
