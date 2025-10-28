package config

// ServerConfig is the config struct for servers (only HTTP_PORT)
type ServerConfig struct {
	HTTPPort int `env:"HTTP_PORT" envDefault:"8080"`
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
	Exchange string `env:"EXCHANGE"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	VHost    string `env:"VHOST"`

	QueueSend string `env:"QUEUE"`
}

/*
 Nah, dude!!! Too much complexity!

// QueueSend is the config that lists MQ queues to send notifications into (by channel)
type QueueSend struct {
	Email    string `env:"EMAIL"`
	Telegram string `env:"TELEGRAM"`
	Console  string `env:"CONSOLE"`
}

*/

// RetryStrategyConfig is the retry strategy config struct
//
// specifies how retry operations will be handled
//
// supposed to be used for multiple things like RABBITMQ_RETRIES, EMAIL_RETRIES, etc.
type RetryStrategyConfig struct {
	Attempts          int     `env:"ATTEMPTS" envDefault:"3"`
	DelayMilliseconds int     `env:"DELAY_MILLISECONDS" envDefault:"500"`
	Backoff           float64 `env:"BACKOFF" envDefault:"1"`
}

// PostgresConfig is the postgres connections config struct
type PostgresConfig struct {
	MasterDSN                    string   `env:"MASTER_DSN"`
	SlaveDSNs                    []string `env:"SLAVE_DSNS" envSeparator:","`
	MaxOpenConnections           int      `env:"MAX_OPEN_CONNECTIONS" envDefault:"3"`
	MaxIdleConnections           int      `env:"MAX_IDLE_CONNECTIONS" envDefault:"5"`
	ConnectionMaxLifetimeSeconds int      `env:"CONNECTION_MAX_LIFETIME_SECONDS" envDefault:"0"`
}

// RedisConfig is the redis connection config struct
type RedisConfig struct {
	Addr       string `env:"ADDR" envDefault:"localhost:6379"`
	Password   string `env:"PASSWORD" envDefault:""`
	DB         int    `env:"DB" envDefault:"0"`
	TTLSeconds int    `env:"TTL_SECONDS" envDefault:"0"`
}

// FetcherConfig is the config struct for fetching batches periodically
type FetcherConfig struct {
	FetchPeriodSeconds      int `env:"FETCH_PERIOD_SECONDS" envDefault:"60"`
	FetchMaxDiapasonSeconds int `env:"FETCH_MAX_DIAPASON_SECONDS" envDefault:"100"`
}
