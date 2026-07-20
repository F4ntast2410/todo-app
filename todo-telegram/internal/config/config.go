package config

import (
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env                 string `env:"ENV" env-default:"development"`
	BotToken            string `env:"TELEGRAM_BOT_TOKEN"`
	TaskServiceGRPCAddr string `env:"TASK_SERVICE_GRPC_ADDR" env-required:"true"`
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
