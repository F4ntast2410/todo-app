package tgbot

import (
	"context"
	"log/slog"

	"proj/internal/entity"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Определяем интерфейсы бизнес-логики, которые нужны боту
type TaskUsecase interface {
	CreateTask(ctx context.Context, title string, description string, userID int) (*entity.Task, error)
	GetTasksByUserID(ctx context.Context, userID int) ([]entity.Task, error)
}
type UserUsecase interface {
	RegisterUserTg(ctx context.Context, ID int64, username string) error
	GetUserByTgID(ctx context.Context, userID int64, username string) (*entity.UserTg, error)
}

type BotServer struct {
	bot          *tgbotapi.BotAPI
	taskUC       TaskUsecase
	userUC       UserUsecase
	logger       *slog.Logger
	sessionCache *SessionCache
}

func NewBotServer(token string, taskUC TaskUsecase, userUC UserUsecase, logger *slog.Logger) (*BotServer, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &BotServer{
		bot:    bot,
		taskUC: taskUC,
		userUC: userUC,
		logger: logger,
		sessionCache: &SessionCache{
			Cache: make(map[int64]UserSession),
		},
	}, nil
}

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
		if update.CallbackQuery != nil {
			WithLoggingCallback(b.handleCallback)(update.CallbackQuery)
			continue
		}
		if update.Message != nil {
			WithLoggingMessage(b.handleMessage)(update.Message)
			continue
		}
	}
}

func (b *BotServer) Send(msg tgbotapi.Chattable) {
	if _, err := b.bot.Send(msg); err != nil {
		b.logger.Error("failed to send message", slog.String("error", err.Error()))
	}
}

func (b *BotServer) Stop() {
	b.bot.StopReceivingUpdates()
}
