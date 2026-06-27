package main

import (
	"log/slog"
	"net/http"
	"os"
	"proj/internal/config"
	"proj/internal/handler"
	"proj/internal/middleware"
	"proj/internal/repository"
	"proj/internal/tgbot"
	"proj/internal/usecase"
	"sync"

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
	taskUsecase := &usecase.TaskUsecaseImpl{Repo: storage}
	userUsecase := &usecase.UserUsecaseImpl{Repo: storage} // Создай пустую структуру в usecase, если еще нет
	handler := &handler.TaskHandler{UC: taskUsecase, Logger: logger}
	mux := http.NewServeMux()

	// Регистрируем твои роуты на этот mux
	mux.HandleFunc("POST /tasks", handler.CreateTaskHandler)
	// mux.HandleFunc("PUT /tasks/{id}", handler.UpdateTaskHandler)
	// mux.HandleFunc("DELETE /tasks/{id}", handler.DeleteTaskHandler)
	// Допустим, у тебя есть еще GET /tasks
	// mux.HandleFunc("GET /tasks", handler.GetTasksHandler)

	// ОБЕРТЫВАЕМ НАШ РОУТЕР В МИДЛВАРЬ ЛОГИРОВАНИЯ
	wrappedMux := middleware.Logger(mux)

	var wg sync.WaitGroup

	// Указываем, что у нас будет 2 параллельных воркера
	wg.Add(2)

	// 2. Запускаем HTTP-сервер
	go func() {
		defer wg.Done() // Уменьшаем счетчик, когда сервер завершит работу

		logger.Info("starting HTTP server", slog.String("port", cfg.ServerPort))
		err := http.ListenAndServe(":"+cfg.ServerPort, wrappedMux)
		if err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server failed", slog.String("err", err.Error()))
		}
	}()

	// 3. Инициализируем и запускаем бота
	botServer, err := tgbot.NewBotServer(cfg.BotToken, taskUsecase, userUsecase, logger)
	if err != nil {
		logger.Error("failed to initialize bot", slog.String("err", err.Error()))
		os.Exit(1)
	}

	go func() {
		defer wg.Done() // Уменьшаем счетчик, когда бот завершит работу

		logger.Info("starting TG bot")
		// Предполагаем, что твой botServer.Start() блокирует поток.
		// Если внутри него уже есть бесконечный цикл, то этот воркер будет активен, пока бот работает.
		botServer.Start()
	}()

	// 4. Главная блокировка вместо botServer.Start() в основном потоке
	// main() будет послушно ждать здесь, пока оба счетчика (wg) не станут равны 0
	wg.Wait()
	logger.Info("all services stopped")
}
