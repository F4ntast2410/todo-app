package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"proj/internal/client/authclient"
	"proj/internal/config"
	grpcDelivery "proj/internal/delivery/grpc"
	httpHandler "proj/internal/delivery/http"
	"proj/internal/delivery/http/middleware"
	repository "proj/internal/repository/postgres"
	"proj/internal/usecase"
	"sync"
	"syscall"
	"time"
	"todo-proto/task"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
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

	storage := repository.NewPostgressStorage(db)
	taskUsecase := usecase.NewTaskUsecase(storage)
	taskHandler := httpHandler.NewTaskHandler(taskUsecase, logger)
	taskGRPCServer := grpcDelivery.NewTaskServer(taskUsecase)
	grpcServer := grpc.NewServer()
	task.RegisterTaskServiceServer(grpcServer, taskGRPCServer)

	authClient, err := authclient.New(cfg.UserServiceGRPCAddr)
	if err != nil {
		logger.Error("failed to connect to user service", slog.String("error", err.Error()))
		panic(err)
	}
	defer authClient.Close()

	authMiddleware := middleware.Auth(authClient)
	mux := http.NewServeMux()

	mux.Handle("POST /tasks", authMiddleware(http.HandlerFunc(taskHandler.CreateTaskHandler)))
	mux.Handle("PUT /tasks/{id}", authMiddleware(http.HandlerFunc(taskHandler.UpdateTaskHandler)))
	mux.Handle("DELETE /tasks/{id}", authMiddleware(http.HandlerFunc(taskHandler.DeleteTaskHandler)))
	mux.Handle("PUT /tasks/restore/{id}", authMiddleware(http.HandlerFunc(taskHandler.RestoreTaskHandler)))
	mux.Handle("GET /tasks/", authMiddleware(http.HandlerFunc(taskHandler.TaskListHandler)))
	mux.Handle("GET /tasks/trash", authMiddleware(http.HandlerFunc(taskHandler.TrashListHandler)))

	wrappedMux := middleware.Logger(mux)
	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: wrappedMux,
	}
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		logger.Error("grpc listen failed", slog.String("error", err.Error()))
		panic(err)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("HTTP server forced shutdown", slog.String("err", err.Error()))
	}
	grpcServer.GracefulStop()
	wg.Wait()
	logger.Info("all services stopped")
}
