package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceName string `mapstructure:"SERVICE_NAME"`
	NATSURL     string `mapstructure:"NATS_URL"`
	SMTPHost    string `mapstructure:"SMTP_HOST"`
	SMTPPort    string `mapstructure:"SMTP_PORT"`
	SMTPEmail   string `mapstructure:"SMTP_EMAIL"`
	SMTPPass    string `mapstructure:"SMTP_PASSWORD"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("SERVICE_NAME", "email-service")
	v.SetDefault("NATS_URL", "nats://nats:4222")
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
