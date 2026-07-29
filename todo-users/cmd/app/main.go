package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"proj/internal/config"
	grpcDelivery "proj/internal/delivery/grpc"
	httpHandler "proj/internal/delivery/http"
	"proj/internal/delivery/http/middleware"
	"proj/internal/mailer"
	repository "proj/internal/repository/postgres"
	cache "proj/internal/repository/redis"
	"proj/internal/usecase"
	"sync"
	"syscall"
	"time"
	"todo-proto/auth"
	"todo-proto/usertg"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.MustLoad()

	var db *sqlx.DB
	var err error

	maxRetries := 10
	for i := 1; i <= maxRetries; i++ {
		db, err = sqlx.Connect("pgx", cfg.DatabaseURL)
		if err == nil {
			logger.Info("successfully connected to database")
			break
		}

		logger.Warn("database connection failed, retrying...",
			slog.Int("attempt", i),
			slog.Int("max_attempts", maxRetries),
			slog.String("error", err.Error()),
		)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		logger.Error("failed to connect to database after max retries", slog.String("error", err.Error()))
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
	mailer := mailer.New(mailer.Config{
		ClientID:     cfg.GmailClientID,
		ClientSecret: cfg.GmailClientSecret,
		RefreshToken: cfg.GmailRefreshToken,
		From:         cfg.GmailFrom,
	})
	// Инициализируем слои
	cache := cache.NewRedisRepository(client)
	storage := repository.NewPostgressStorage(db)
	userUsecase := usecase.NewUserUsecase(storage, cache, mailer)
	authGRPCServer := grpcDelivery.NewAuthServer(userUsecase)
	grpcServer := grpc.NewServer()
	userTgServer := grpcDelivery.NewUserTgServer(userUsecase) // новый файл, по образцу AuthServer
	usertg.RegisterUserTgServiceServer(grpcServer, userTgServer)
	auth.RegisterAuthServiceServer(grpcServer, authGRPCServer)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		logger.Error("grpc listen failed", slog.String("error", err.Error()))
		panic(err)
	}
	userHandler := httpHandler.NewUserHandler(userUsecase, logger, cfg.BotToken)

	authMiddleware := middleware.Auth(userUsecase)
	mux := http.NewServeMux()

	mux.Handle("GET /auth/me", authMiddleware(http.HandlerFunc(userHandler.GetUserHandler)))

	mux.HandleFunc("POST /auth/register", userHandler.RegisterHandler)
	mux.HandleFunc("POST /auth/register/verify", userHandler.RegisterVerify2FAHandler)
	mux.HandleFunc("POST /auth/login", userHandler.LoginHandler)
	mux.HandleFunc("POST /auth/verify2fa", userHandler.Verify2FA)
	mux.HandleFunc("POST /auth/logout", userHandler.LogoutHandler)

	mux.Handle("PUT /auth/me", authMiddleware(http.HandlerFunc(userHandler.UpdateProfileHandler)))
	mux.Handle("PUT /auth/me/email/verify", authMiddleware(http.HandlerFunc(userHandler.UpdateEmailVerifyHandler)))
	mux.Handle("PUT /auth/me/password", authMiddleware(http.HandlerFunc(userHandler.ChangePasswordHandler)))
	mux.Handle("GET /auth/telegram/link", authMiddleware(http.HandlerFunc(userHandler.GetTelegramLinkHandler)))
	mux.Handle("POST /auth/telegram/link", authMiddleware(http.HandlerFunc(userHandler.LinkTelegramHandler)))
	mux.HandleFunc("POST /auth/telegram/login", userHandler.TelegramLoginHandler)

	mux.Handle("POST /auth/me/link/email", authMiddleware(http.HandlerFunc(userHandler.SetEmailHandler)))
	mux.Handle("POST /auth/me/link/verify2fa", authMiddleware(http.HandlerFunc(userHandler.EmailSetVerify2FAHandler)))
	mux.Handle("POST /auth/me/link/password", authMiddleware(http.HandlerFunc(userHandler.SetEmailPasswordHandler)))

	wrappedMux := middleware.Logger(mux)
	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: wrappedMux,
	}

	// Канал для аварийных ошибок — буфер 2, чтобы обе горутины могли записать не блокируясь
	fatalErr := make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	// Запускаем HTTP-сервер
	go func() {
		defer wg.Done()
		logger.Info("starting HTTP server", slog.String("port", cfg.ServerPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server failed", slog.String("err", err.Error()))
			fatalErr <- err // сообщаем об аварии
		}
	}()
	go func() {
		defer wg.Done()
		logger.Info("starting gRPC server", slog.String("port", cfg.GRPCPort))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("gRPC server failed", slog.String("err", err.Error()))
			fatalErr <- err
		}
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

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("HTTP server forced shutdown", slog.String("err", err.Error()))
	}
	grpcServer.GracefulStop()
	wg.Wait()
	logger.Info("all services stopped")
}
