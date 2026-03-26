package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/HungryArthur/go-shortener/internal/config"
	"github.com/HungryArthur/go-shortener/internal/db"
	"github.com/HungryArthur/go-shortener/internal/handlers/ping"
	url_handler "github.com/HungryArthur/go-shortener/internal/handlers/url"
	"github.com/HungryArthur/go-shortener/internal/middlewares"
	repository "github.com/HungryArthur/go-shortener/internal/repository/url"
	"github.com/HungryArthur/go-shortener/internal/service"
)

func main() {
	config.Load()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	sugar.Infow("config load",
		"run address", config.FlagRunAddr,
		"run address", config.FileStoragePath,
		"run address", maskDSN(config.FileDatabaseDSN),
	)

	var repo service.URLRepository

	if config.FileDatabaseDSN != "" {
		sugar.Info("using PostgreSQL repository")

		// Подключаемся к БД
		database, err := db.Connect(config.FileDatabaseDSN, logger)
		if err != nil {
			sugar.Fatalw("failed to connect to database", "error", err)
		}
		defer database.Close()

		repo = repository.NewURLPostgresRepository(database, logger)
	} else if config.FileStoragePath != "" {

		sugar.Infow("using file storage", "path", config.FileStoragePath)
		repo = repository.NewURLDiskJSONRepository(logger, config.FileStoragePath)

	} else {
		sugar.Info("using in-memory storage")
		repo = repository.NewURLMemoryRepository()
	}

	urlService := service.NewURLService(repo)
	urlHandler := url_handler.NewURLHandler(urlService)
	pingHandler := ping.NewPingHandler(repo, logger)

	router := chi.NewRouter()

	router.Use(middlewares.WriteHeader(logger))
	router.Use(middlewares.GzipMiddleware)

	router.Get("/ping", pingHandler.Ping)
	router.Get("/{shortenedURL}", urlHandler.GetTextPlain)
	router.Post("/", urlHandler.CreateTextPlain)
	router.Post("/api/shorten", urlHandler.CreateJSON)

	srv := http.Server{Addr: config.FlagRunAddr, Handler: router}

	sugar.Infow("starting server", "address", config.FlagRunAddr)

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			sugar.Errorw("can't start server", err)
		}
	}()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	quitCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = srv.Shutdown(quitCtx)
	if err != nil {
		sugar.Error("can't shutdown http server")
	} else {
		sugar.Info("successfully shutdowned http server")
	}

}

func maskDSN(dsn string) string {
	if dsn == "" {
		return ""
	}
	return "postgres://****:****@****/****"
}
