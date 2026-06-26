package tgbot

import (
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotServer) Start() {
	b.logger.Info("Telegram bot started", slog.String("botname", b.bot.Self.UserName))

	// Настраиваем конфигурацию получения обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Получаем канал с обновлениями от Telegram
	updates := b.bot.GetUpdatesChan(u)

	// Читаем сообщения из канала в бесконечном цикле
	for update := range updates {
		// Игнорируем любые обновления, кроме текстовых сообщений
		if update.Message == nil {
			continue
		}

		// Обрабатываем сообщение
		WithLogging(b.handleMessage)(update.Message)
	}
}
