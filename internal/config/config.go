package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `env:"ENV" env-default:"development"`
	ServerPort  string `env:"SERVER_PORT" env-default:"8080"`
	DatabaseURL string `env:"DATABASE_URL" env-required:"true"`
	BotToken    string `env:"TELEGRAM_BOT_TOKEN"`
}

// MustLoad считывает конфиг один раз при старте.
// Префикс Must в Go означает, что если функция не сработает — приложение упадет (panic).
func MustLoad() *Config {
	var cfg Config

	// Читаем файл .env, если он есть рядом
	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		// Если файла .env нет (например, в Docker-контейнере на проде),
		// библиотека попытается прочитать переменные напрямую из системы
		err = cleanenv.ReadEnv(&cfg)
		if err != nil {
			log.Fatalf("не удалось прочитать конфигурацию: %v", err)
		}
	}

	return &cfg
}
