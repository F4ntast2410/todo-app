package config

import (
	"fmt"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string `env:"ENV" env-default:"development"`
	ServerPort   string `env:"SERVER_PORT" env-default:"8080"`
	GRPCPort     string `env:"GRPC_PORT" env-default:"9090"`
	DatabaseURL  string `env:"DATABASE_URL" env-required:"true"`
	BotToken     string `env:"TELEGRAM_BOT_TOKEN"`
	Redis        RedisConfig
	SMTPHost     string `env:"SMTP_HOST" env-default:"smtp.gmail.com"`
	SMTPPort     string `env:"SMTP_PORT" env-default:"587"`
	SMTPUsername string `env:"SMTP_USERNAME" env-required:"true"`
	SMTPPassword string `env:"SMTP_PASSWORD" env-required:"true"`
	SMTPFrom     string `env:"SMTP_FROM" env-required:"true"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" env-default:"localhost"`
	Port     string `env:"REDIS_PORT" env-default:"6379"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB" env-default:"0"`
}

// Addr собирает host:port для redis.Options — так проще, чем городить URL
func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

func MustLoad() *Config {
	var cfg Config

	// 1. Сначала читаем .env (если он есть)
	_ = cleanenv.ReadConfig(".env", &cfg)

	// 2. Поверх перезаписываем реальными переменными окружения из Docker/OS
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		slog.Error("failed to load env config", slog.String("error", err.Error()))
		panic(err)
	}

	return &cfg
}
