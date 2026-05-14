package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv      string `mapstructure:"APP_ENV"`
	ServiceName string `mapstructure:"SERVICE_NAME"`

	GRPCPort    string `mapstructure:"GRPC_PORT"`
	HTTPPort    string `mapstructure:"HTTP_PORT"`
	MetricsPort string `mapstructure:"METRICS_PORT"`

	PostgresDSN string `mapstructure:"POSTGRES_DSN"`

	Redis RedisConfig `mapstructure:",squash"`
	NATS  NATSConfig  `mapstructure:",squash"`

	CacheTTL time.Duration `mapstructure:"CACHE_TTL"`

	JWTSecret string `mapstructure:"JWT_SECRET"`

	OTELEnabled  bool   `mapstructure:"OTEL_ENABLED"`
	OTELService  string `mapstructure:"OTEL_SERVICE"`
	OTELEndpoint string `mapstructure:"OTEL_ENDPOINT"`

	SMTP SMTPConfig `mapstructure:",squash"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"REDIS_ADDR"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

type NATSConfig struct {
	URL string `mapstructure:"NATS_URL"`
}

type SMTPConfig struct {
	Host     string `mapstructure:"SMTP_HOST"`
	Port     string `mapstructure:"SMTP_PORT"`
	Email    string `mapstructure:"SMTP_EMAIL"`
	Password string `mapstructure:"SMTP_PASSWORD"`
}

func Load() (*Config, error) {

	v := viper.New()

	v.SetConfigFile(".env")
	v.SetConfigType("env")

	v.SetEnvKeyReplacer(
		strings.NewReplacer(".", "_"),
	)

	v.AutomaticEnv()

	// App
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("SERVICE_NAME", "user-service")

	// Ports
	v.SetDefault("GRPC_PORT", "50052")
	v.SetDefault("HTTP_PORT", "8082")
	v.SetDefault("METRICS_PORT", "9102")

	// Database
	v.SetDefault(
		"POSTGRES_DSN",
		"postgres://postgres:0000@localhost:5434/userdb?sslmode=disable",
	)

	// Redis
	v.SetDefault("REDIS_ADDR", "localhost:6380")
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)

	// NATS
	v.SetDefault(
		"NATS_URL",
		"nats://localhost:4223",
	)

	// Cache
	v.SetDefault("CACHE_TTL", "5m")

	// JWT
	v.SetDefault("JWT_SECRET", "super-secret-key")

	// OpenTelemetry
	v.SetDefault("OTEL_ENABLED", false)
	v.SetDefault("OTEL_SERVICE", "user-service")

	v.SetDefault(
		"OTEL_ENDPOINT",
		"localhost:4317",
	)

	// SMTP
	v.SetDefault("SMTP_HOST", "smtp.gmail.com")
	v.SetDefault("SMTP_PORT", "587")
	v.SetDefault("SMTP_EMAIL", "")
	v.SetDefault("SMTP_PASSWORD", "")

	_ = v.ReadInConfig()

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
