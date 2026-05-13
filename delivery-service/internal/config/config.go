package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv      string               `mapstructure:"APP_ENV"`
	ServiceName string               `mapstructure:"SERVICE_NAME"`
	GRPCPort    string               `mapstructure:"GRPC_PORT"`
	MetricsPort string               `mapstructure:"METRICS_PORT"`
	PostgresDSN string               `mapstructure:"POSTGRES_DSN"`
	Redis       RedisConfig          `mapstructure:",squash"`
	NATS        NATSConfig           `mapstructure:",squash"`
	Restaurant  RestaurantGRPCConfig `mapstructure:",squash"`
	CacheTTL    time.Duration        `mapstructure:"CACHE_TTL"`
	ETACacheTTL time.Duration        `mapstructure:"ETA_CACHE_TTL"`
	JWTSecret   string               `mapstructure:"JWT_SECRET"`
	OTELEnabled bool                 `mapstructure:"OTEL_ENABLED"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"REDIS_ADDR"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

type NATSConfig struct {
	URL string `mapstructure:"NATS_URL"`
}

type RestaurantGRPCConfig struct {
	Addr    string        `mapstructure:"RESTAURANT_GRPC_ADDR"`
	Timeout time.Duration `mapstructure:"RESTAURANT_TIMEOUT"`
	Retries int           `mapstructure:"RESTAURANT_RETRIES"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("APP_ENV", "development")
	v.SetDefault("SERVICE_NAME", "delivery-service")
	v.SetDefault("GRPC_PORT", "50056")
	v.SetDefault("METRICS_PORT", "9106")
	v.SetDefault("CACHE_TTL", "5m")
	v.SetDefault("ETA_CACHE_TTL", "2m")
	v.SetDefault("REDIS_ADDR", "localhost:6379")
	v.SetDefault("NATS_URL", "nats://localhost:4222")
	v.SetDefault("RESTAURANT_GRPC_ADDR", "localhost:50055")
	v.SetDefault("RESTAURANT_TIMEOUT", "2s")
	v.SetDefault("RESTAURANT_RETRIES", 2)

	_ = v.ReadInConfig()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
