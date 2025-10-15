package config

import (
	"fmt"
	"github.com/wb-go/wbf/config"
	"strings"
)

// AppConfig is THE whole config struct
type AppConfig struct {
	ServerConfig ServerConfig `env-prefix:"SERVER_"`
	LogConfig    LogConfig    `env-prefix:"LOG_"`

	RabbitMQConfig RabbitMQConfig `env-prefix:"RABBITMQ_"`

	PostgresConfig PostgresConfig `env-prefix:"POSTGRES_"`
	RedisConfig    RedisConfig    `env-prefix:"REDIS_"`

	PostgresRetryConfig RetryStrategyConfig `env-prefix:"RETRY_POSTGRES_"`
	RabbitMQRetryConfig RetryStrategyConfig `env-prefix:"RETRY_RABBITMQ_"`
	RedisRetryConfig    RetryStrategyConfig `env-prefix:"RETRY_REDIS_"`

	FetcherConfig FetcherConfig `env-prefix:"FETCHER_"`
}

// NewAppConfig creates a new struct of "THE config"
//
// has some defaults
func NewAppConfig(configFilePath, envFilePath string) (*AppConfig, error) {
	appConfig := &AppConfig{}

	cfg := config.New()

	//region defaults
	cfg.SetDefault("delayed_notifier.server.http.port", 8080)
	cfg.SetDefault("delayed_notifier.log.level", "info")

	cfg.SetDefault("delayed_notifier.postgres.max_open_connections", 2)
	cfg.SetDefault("delayed_notifier.postgres.max_idle_connections", 2)
	cfg.SetDefault("delayed_notifier.postgres.connection_max_lifetime_seconds", 0)

	cfg.SetDefault("delayed_notifier.redis.db", 0)
	cfg.SetDefault("delayed_notifier.redis.ttl_seconds", 20)

	cfg.SetDefault("delayed_notifier.retry_redis.attempts", 3)
	cfg.SetDefault("delayed_notifier.retry_postgres.attempts", 3)
	cfg.SetDefault("delayed_notifier.retry_rabbitmq.attempts", 3)

	cfg.SetDefault("delayed_notifier.retry_redis.delay_milliseconds", 300)
	cfg.SetDefault("delayed_notifier.retry_postgres.delay_milliseconds", 300)
	cfg.SetDefault("delayed_notifier.retry_rabbitmq.delay_milliseconds", 300)

	cfg.SetDefault("delayed_notifier.retry_redis.backoff", 1.5)
	cfg.SetDefault("delayed_notifier.retry_postgres.backoff", 1.5)
	cfg.SetDefault("delayed_notifier.retry_rabbitmq.backoff", 1.5)
	//endregion

	// region flags

	// why flags lol dude

	// does it work as "set as default"? if I specify 8080 as default flag value, then it's always non-empty
	_ = cfg.DefineFlag("p", "http_port", "SERVER_HTTP_PORT", 8080, "HTTP server port")

	_ = cfg.DefineFlag("l", "log_level", "LOG_LEVEL", "info", "Log level (lowercase)")

	// how about that? If it's empty, will it override my .env?
	// _ = cfg.DefineFlag("q", "queue", "RABBITMQ_QUEUE", "", "RabbitMQ queue name")

	// cfg.ParseFlags()
	//endregion

	err := cfg.Load(configFilePath, envFilePath, "")
	if err != nil {
		return appConfig, fmt.Errorf("failed to load config: %w", err)
	}

	// 1. ServerConfig
	appConfig.ServerConfig.HTTPPort = cfg.GetInt("delayed_notifier.server.http.port")

	// 2. LogConfig
	appConfig.LogConfig.LogLevel = cfg.GetString("delayed_notifier.log.level")

	// 3. RabbitMQConfig
	appConfig.RabbitMQConfig.Exchange = cfg.GetString("delayed_notifier.rabbitmq.exchange")
	appConfig.RabbitMQConfig.User = cfg.GetString("delayed_notifier.rabbitmq.user")
	appConfig.RabbitMQConfig.Password = cfg.GetString("delayed_notifier.rabbitmq.password")
	appConfig.RabbitMQConfig.Host = cfg.GetString("delayed_notifier.rabbitmq.host")
	appConfig.RabbitMQConfig.Port = cfg.GetInt("delayed_notifier.rabbitmq.port")
	appConfig.RabbitMQConfig.VHost = cfg.GetString("delayed_notifier.rabbitmq.vhost")
	appConfig.RabbitMQConfig.QueueSend.Email = cfg.GetString("delayed_notifier.rabbitmq.queue_send.email")
	appConfig.RabbitMQConfig.QueueSend.Telegram = cfg.GetString("delayed_notifier.rabbitmq.queue_send.telegram")
	appConfig.RabbitMQConfig.QueueSend.Console = cfg.GetString("delayed_notifier.rabbitmq.queue_send.console")

	// 4. PostgresConfig
	appConfig.PostgresConfig.MasterDSN = cfg.GetString("delayed_notifier.postgres.master_dsn")
	appConfig.PostgresConfig.SlaveDSNs = strings.Split(cfg.GetString("delayed_notifier.postgres.slave_dsns"), ",")
	appConfig.PostgresConfig.MaxOpenConnections = cfg.GetInt("delayed_notifier.postgres.max_open_connections")
	appConfig.PostgresConfig.MaxIdleConnections = cfg.GetInt("delayed_notifier.postgres.max_idle_connections")
	appConfig.PostgresConfig.ConnectionMaxLifetimeSeconds = cfg.GetInt("delayed_notifier.postgres.connection_max_lifetime_seconds")

	// 5. RedisConfig
	appConfig.RedisConfig.Addr = cfg.GetString("delayed_notifier.redis.addr")
	appConfig.RedisConfig.Password = cfg.GetString("delayed_notifier.redis.password")
	appConfig.RedisConfig.DB = cfg.GetInt("delayed_notifier.redis.db")
	appConfig.RedisConfig.TTLSeconds = cfg.GetInt("delayed_notifier.redis.ttl_seconds")

	//region 6-8. Retries
	appConfig.PostgresRetryConfig.Attempts = cfg.GetInt("delayed_notifier.retry_postgres.attempts")
	appConfig.PostgresRetryConfig.DelayMilliseconds = cfg.GetInt("delayed_notifier.retry_postgres.delay_milliseconds")
	appConfig.PostgresRetryConfig.Backoff = cfg.GetFloat64("delayed_notifier.retry_postgres.backoff")

	appConfig.RabbitMQRetryConfig.Attempts = cfg.GetInt("delayed_notifier.retry_rabbitmq.attempts")
	appConfig.RabbitMQRetryConfig.DelayMilliseconds = cfg.GetInt("delayed_notifier.retry_rabbitmq.delay_milliseconds")
	appConfig.RabbitMQRetryConfig.Backoff = cfg.GetFloat64("delayed_notifier.retry_rabbitmq.backoff")

	appConfig.RedisRetryConfig.Attempts = cfg.GetInt("delayed_notifier.retry_redis.attempts")
	appConfig.RedisRetryConfig.DelayMilliseconds = cfg.GetInt("delayed_notifier.retry_redis.delay_milliseconds")
	appConfig.RedisRetryConfig.Backoff = cfg.GetFloat64("delayed_notifier.retry_redis.backoff")
	//endregion

	//9. FetcherConfig
	appConfig.FetcherConfig.FetchPeriodSeconds = cfg.GetInt("delayed_notifier.fetcher.fetch_period_seconds")

	return appConfig, nil
}
