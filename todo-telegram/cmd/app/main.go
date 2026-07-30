package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"proj/internal/client/taskclient"
	"proj/internal/client/usertgclient"
	"proj/internal/config"
	"proj/internal/delivery/tgbot"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 🟢 ПРИНУДИТЕЛЬНО ОТПРАВЛЯЕМ ВСЕ HTTP-ЗАПРОСЫ ПО IPv6 (tcp6)
	http.DefaultTransport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "tcp6", addr)
		},
	}

	cfg := config.MustLoad()

	taskClient, err := taskclient.New(cfg.TaskServiceGRPCAddr)
	if err != nil {
		logger.Error("failed to connect to task service", slog.String("error", err.Error()))
		panic(err)
	}
	defer taskClient.Close()

	userTgClient, err := usertgclient.New(cfg.UserServiceGRPCAddr)
	if err != nil {
		logger.Error("failed to connect to usertg service", slog.String("error", err.Error()))
		panic(err)
	}
	defer userTgClient.Close()

	botServer, err := tgbot.NewBotServer(cfg.BotToken, taskClient, userTgClient, logger)
	if err != nil {
		logger.Error("failed to initialize bot", slog.String("err", err.Error()))
		os.Exit(1)
	}

	// Запускаем бота
	go func() {
		logger.Info("starting TG bot")
		botServer.Start()
	}()

	// Ждём сигнала остановки (Ctrl+C или docker stop)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	logger.Info("shutdown signal received", slog.String("signal", sig.String()))

	// Дальше — тот же код остановки
	logger.Info("stopping services...")
	botServer.Stop()
	logger.Info("all services stopped")
}
