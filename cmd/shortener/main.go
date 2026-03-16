package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HungryArthur/go-shortener/internal/config"
	url_handler "github.com/HungryArthur/go-shortener/internal/handlers/url"
	"github.com/HungryArthur/go-shortener/internal/middlewares"
	"github.com/HungryArthur/go-shortener/internal/repository"
	"github.com/HungryArthur/go-shortener/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	config.Load()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	
	sugar := logger.Sugar()
	
	config.Load()
	sugar.Infow("config load", "run address", config.FlagRunAddr)

	repo := repository.NewURLRepository()
	service := service.NewURLService(repo)
	handler := url_handler.NewURLHandler(service)

	router := chi.NewRouter()
	
	router.Use(middlewares.WriteHeader(logger))

	router.Get("/{shortenedURL}", handler.GetTextPlain)
	router.Post("/", handler.CreateTextPlain)

	router.Post("/api/shorten", handler.GetJson)

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
