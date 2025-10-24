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

	UniversalQueue string `env:"QUEUE"`
	// QueueForChannel QueueRead `env:"QUEUE_"`
	Consumer string `env:"CONSUMER"`
	AutoAck  bool   `env:"AUTO_ACK" envDefault:"false"`
	NoWait   bool   `env:"NO_WAIT" envDefault:"false"`
}

// QueueRead is the config that lists MQ queues to read notifications from (by channel)
type QueueRead struct {
	Email    string `env:"EMAIL"`
	Telegram string `env:"TELEGRAM"`
	Console  string `env:"CONSOLE"`
}

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
