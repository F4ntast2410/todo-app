package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"proj/internal/config"
	httpHandler "proj/internal/delivery/http"
	"proj/internal/delivery/tgbot"
	"proj/internal/middleware"
	"proj/internal/repository"
	"proj/internal/usecase"
	"sync"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.MustLoad()

	db, err := sqlx.Connect("pgx", cfg.DatabaseURL)
	if err != nil {
		logger.Error("database connection failed", slog.String("error", err.Error()))
		panic(err)
	}
	defer db.Close()

	// Инициализируем слои
	storage := &repository.PostgresStorage{DB: db}
	taskUsecase := &usecase.TaskUsecaseImpl{TaskRepo: storage}
	userUsecase := &usecase.UserUsecaseImpl{UserRepo: storage} // Создай пустую структуру в usecase, если еще нет
	handler := &httpHandler.TaskHandler{UC: taskUsecase, Logger: logger}
	mux := http.NewServeMux()

	// Регистрируем твои роуты на этот mux
	mux.HandleFunc("POST /tasks", handler.CreateTaskHandler)
	// mux.HandleFunc("PUT /tasks/{id}", handler.UpdateTaskHandler)
	// mux.HandleFunc("DELETE /tasks/{id}", handler.DeleteTaskHandler)
	// Допустим, у тебя есть еще GET /tasks
	// mux.HandleFunc("GET /tasks", handler.GetTasksHandler)

	// ОБЕРТЫВАЕМ НАШ РОУТЕР В МИДЛВАРЬ ЛОГИРОВАНИЯ
	wrappedMux := middleware.Logger(mux)

	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: wrappedMux,
	}

	botServer, err := tgbot.NewBotServer(cfg.BotToken, taskUsecase, userUsecase, logger)
	if err != nil {
		logger.Error("failed to initialize bot", slog.String("err", err.Error()))
		os.Exit(1)
	}

	var wg sync.WaitGroup

	wg.Add(2)

	// Запускаем HTTP-сервер
	go func() {
		defer wg.Done()
		logger.Info("starting HTTP server", slog.String("port", cfg.ServerPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server failed", slog.String("err", err.Error()))
		}
	}()

	// Запускаем бота
	go func() {
		defer wg.Done()
		logger.Info("starting TG bot")
		botServer.Start() // блокируется до вызова Stop()
	}()

	// Ждём сигнала остановки (Ctrl+C или docker stop)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutdown signal received, stopping services...")

	// Останавливаем бота
	botServer.Stop()

	// Останавливаем HTTP-сервер — даём 5 секунд на завершение активных запросов
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("HTTP server forced shutdown", slog.String("err", err.Error()))
	}

	// Ждём пока оба горутина завершатся
	wg.Wait()
	logger.Info("all services stopped")
}
