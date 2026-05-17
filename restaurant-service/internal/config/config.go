package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv      string        `mapstructure:"APP_ENV"`
	ServiceName string        `mapstructure:"SERVICE_NAME"`
	GRPCPort    string        `mapstructure:"GRPC_PORT"`
	MetricsPort string        `mapstructure:"METRICS_PORT"`
	PostgresDSN string        `mapstructure:"POSTGRES_DSN"`
	Redis       RedisConfig   `mapstructure:",squash"`
	NATS        NATSConfig    `mapstructure:",squash"`
	CacheTTL    time.Duration `mapstructure:"CACHE_TTL"`
	JWTSecret   string        `mapstructure:"JWT_SECRET"`
	OTELEnabled bool          `mapstructure:"OTEL_ENABLED"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"REDIS_ADDR"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

type NATSConfig struct {
	URL string `mapstructure:"NATS_URL"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("APP_ENV", "development")
	v.SetDefault("SERVICE_NAME", "restaurant-service")
	v.SetDefault("GRPC_PORT", "50055")
	v.SetDefault("METRICS_PORT", "9105")
	v.SetDefault(
		"POSTGRES_DSN",
		"postgres://restaurant:restaurant@restaurant-db:5432/restaurant_service?sslmode=disable",
	)
	v.SetDefault("CACHE_TTL", "5m")
	v.SetDefault("REDIS_ADDR", "localhost:6379")
	v.SetDefault("NATS_URL", "nats://localhost:4222")

	_ = v.ReadInConfig()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
