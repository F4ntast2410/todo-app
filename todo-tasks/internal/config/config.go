package config

import (
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env                 string `env:"ENV" env-default:"development"`
	ServerPort          string `env:"SERVER_PORT" env-default:"8080"`
	GRPCPort            string `env:"GRPC_PORT" env-default:"9090"`
	DatabaseURL         string `env:"DATABASE_URL" env-required:"true"`
	UserServiceGRPCAddr string `env:"USER_SERVICE_GRPC_ADDR" env-required:"true"`
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
