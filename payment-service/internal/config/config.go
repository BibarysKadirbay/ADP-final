package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceName string `mapstructure:"SERVICE_NAME"`
	GRPCPort    string `mapstructure:"GRPC_PORT"`
	MetricsPort string `mapstructure:"METRICS_PORT"`
	PostgresDSN string `mapstructure:"POSTGRES_DSN"`
	NATSURL     string `mapstructure:"NATS_URL"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("SERVICE_NAME", "payment-service")
	v.SetDefault("GRPC_PORT", "50053")
	v.SetDefault("METRICS_PORT", "9103")
	v.SetDefault("POSTGRES_DSN", "postgres://postgres:0000@payment-db:5432/paymentdb?sslmode=disable")
	v.SetDefault("NATS_URL", "nats://nats:4222")

	_ = v.ReadInConfig()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
