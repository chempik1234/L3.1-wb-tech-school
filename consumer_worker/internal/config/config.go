package config

import (
	"fmt"
	"github.com/wb-go/wbf/config"
)

// AppConfig is THE whole config struct
type AppConfig struct {
	LogConfig      LogConfig      `env-prefix:"LOG_"`
	RabbitMQConfig RabbitMQConfig `env-prefix:"RABBITMQ_"`
}

// LogConfig is the config struct for logging
//
// available log levels: "trace", "debug", "info", "warn", "error", "fatal", "panic"
type LogConfig struct {
	LogLevel string `env:"LEVEL" envDefault:"info"`
}

// RabbitMQConfig is the config struct for rabbitMQ (consumer)
//
// only the stated props here are supposed to be changeable
type RabbitMQConfig struct {
	Queue    string `env:"QUEUE"`
	Consumer string `env:"CONSUMER"`
	AutoAck  bool   `env:"AUTO_ACK" envDefault:"false"`
	NoWait   bool   `env:"NO_WAIT" envDefault:"false"`
}

// RetryStrategyConfig is the retry strategy config struct
//
// specifies how retry operations will be handled
//
// supposed to be used for multiple things like RABBITMQ_RETRIES, EMAIL_RETRIES, etc.
type RetryStrategyConfig struct {
	Attempts          int `env:"ATTEMPTS" envDefault:"3"`
	DelayMilliseconds int `env:"DELAY_MILLISECONDS" envDefault:"500"`
	Backoff           int `env:"BACKOFF" envDefault:"1"`
}

// NewAppConfig creates a new struct of "THE config"
func NewAppConfig(configFilePath, envFilePath string) (*AppConfig, error) {
	appConfig := &AppConfig{}

	cfg := config.New()
	cfg.SetDefault("LOG_LEVEL", "info")
	cfg.SetDefault("RABBITMQ_NO_WAIT", false)
	cfg.SetDefault("RABBITMQ_AUTO_ACK", false)

	// region flags

	// why flags lol dude

	// does it work as "set as default"? if I specify 8080 as default flag value, then it's always non-empty
	_ = cfg.DefineFlag("p", "http_port", "SERVER_HTTP_PORT", 8080, "HTTP server port")

	_ = cfg.DefineFlag("l", "log_level", "LOG_LEVEL", "info", "Log level (lowercase)")

	// how about that? If it's empty, will it override my .env?
	// _ = cfg.DefineFlag("q", "queue", "RABBITMQ_QUEUE", "", "RabbitMQ queue name")

	cfg.ParseFlags()
	//endregion

	err := cfg.Load(configFilePath, envFilePath, "CONSUMER_WORKER_")
	if err != nil {
		return appConfig, fmt.Errorf("failed to load config: %w", err)
	}

	// LogConfig
	appConfig.LogConfig.LogLevel = cfg.GetString("LOG_LEVEL")

	// RabbitMQConfig
	appConfig.RabbitMQConfig.Queue = cfg.GetString("RABBITMQ_QUEUE")
	appConfig.RabbitMQConfig.Consumer = cfg.GetString("RABBITMQ_CONSUMER")
	appConfig.RabbitMQConfig.AutoAck = cfg.GetBool("RABBITMQ_AUTO_ACK")
	appConfig.RabbitMQConfig.NoWait = cfg.GetBool("RABBITMQ_NO_WAIT")

	return appConfig, nil
}
