package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"proj/internal/config"
	httpHandler "proj/internal/delivery/http"
	"proj/internal/delivery/http/middleware"
	"proj/internal/delivery/tgbot"
	repository "proj/internal/repository/postgres"
	cache "proj/internal/repository/redis"
	"proj/internal/usecase"
	"sync"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
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

	opt := &redis.Options{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}
	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error("redis connection failed", slog.String("error", err.Error()))
		panic(err)
	}

	// Инициализируем слои
	cache := cache.NewRedisRepository(client)
	storage := repository.NewPostgressStorage(db)
	taskUsecase := usecase.NewTaskUsecase(storage)
	userUsecase := usecase.NewUserUsecase(storage, cache)
	taskHandler := httpHandler.NewTaskHandler(taskUsecase, logger)
	userHandler := httpHandler.NewUserHandler(userUsecase, logger)

	authMiddleware := middleware.Auth(userUsecase)
	mux := http.NewServeMux()

	mux.Handle("POST /tasks", authMiddleware(http.HandlerFunc(taskHandler.CreateTaskHandler)))
	mux.Handle("PUT /tasks/{id}", authMiddleware(http.HandlerFunc(taskHandler.UpdateTaskHandler)))
	mux.Handle("DELETE /tasks/{id}", authMiddleware(http.HandlerFunc(taskHandler.DeleteTaskHandler)))
	mux.Handle("PUT /tasks/restore/{id}", authMiddleware(http.HandlerFunc(taskHandler.RestoreTaskHandler)))
	mux.Handle("GET /tasks/", authMiddleware(http.HandlerFunc(taskHandler.TaskListHandler)))
	mux.Handle("GET /tasks/trash", authMiddleware(http.HandlerFunc(taskHandler.TrashListHandler)))

	mux.Handle("GET /auth/me", authMiddleware(http.HandlerFunc(userHandler.GetUserHandler)))

	mux.HandleFunc("POST /auth/register", userHandler.RegisterHandler)
	mux.HandleFunc("POST /auth/login", userHandler.LoginHandler)
	mux.HandleFunc("POST /auth/verify2fa", userHandler.Verify2FA)
	mux.HandleFunc("POST /auth/logout", userHandler.LogoutHandler)

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

	// Канал для аварийных ошибок — буфер 2, чтобы обе горутины могли записать не блокируясь
	fatalErr := make(chan error, 2)

	// Запускаем HTTP-сервер
	go func() {
		defer wg.Done()
		logger.Info("starting HTTP server", slog.String("port", cfg.ServerPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server failed", slog.String("err", err.Error()))
			fatalErr <- err // сообщаем об аварии
		}
	}()

	// Запускаем бота
	go func() {
		defer wg.Done()
		logger.Info("starting TG bot")
		botServer.Start()
	}()

	// Ждём сигнала остановки (Ctrl+C или docker stop)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Ждём: либо сигнал от ОС, либо авария одного из сервисов
	select {
	case sig := <-quit:
		logger.Info("shutdown signal received", slog.String("signal", sig.String()))
	case err := <-fatalErr:
		logger.Error("fatal error, initiating shutdown", slog.String("err", err.Error()))
	}

	// Дальше — тот же код остановки
	logger.Info("stopping services...")
	botServer.Stop()

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("HTTP server forced shutdown", slog.String("err", err.Error()))
	}

	wg.Wait()
	logger.Info("all services stopped")
}
