package tgbot

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"proj/internal/entity"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Определяем интерфейсы бизнес-логики, которые нужны боту
type TaskUsecase interface {
	CreateTask(ctx context.Context, title string, description string, userID int) (*entity.Task, error)
	GetTasksByUserID(ctx context.Context, userID int) ([]entity.Task, error)
	GetRemovedTasksByUserID(ctx context.Context, userID int) ([]entity.Task, error)
	GetTask(ctx context.Context, taskID int) (*entity.Task, error)
	DeleteTask(ctx context.Context, id int) error
	DeleteForeverTask(ctx context.Context, id int) error
	RecoverTask(ctx context.Context, id int) error
	UpdateDescription(ctx context.Context, taskID int, newDesc string) error
	MarkAsDone(ctx context.Context, id int, status bool) error
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

// newIPv4OnlyHTTPClient возвращает http.Client, который всегда открывает
// соединение по tcp4, независимо от того, что вернёт DNS-резолвер (A/AAAA)
// и что запросит вызывающий код. Нужен, чтобы бот не пытался ходить
// к api.telegram.org по IPv6 внутри Docker-сети, где IPv6-маршрута нет.
func newIPv4OnlyHTTPClient() *http.Client {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, "tcp4", addr)
		},
	}
	return &http.Client{Transport: transport}
}

func NewBotServer(token string, taskUC TaskUsecase, userUC UserUsecase, logger *slog.Logger) (*BotServer, error) {
	bot, err := tgbotapi.NewBotAPIWithClient(token, tgbotapi.APIEndpoint, newIPv4OnlyHTTPClient())
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
		if strings.Contains(err.Error(), "message is not modified") {
			return // ничего страшного, просто пропускаем
		}
		b.logger.Error("failed to send message", slog.String("error", err.Error()))
	}
}

func (b *BotServer) Stop() {
	b.bot.StopReceivingUpdates()
}
