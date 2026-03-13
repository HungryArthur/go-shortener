package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HungryArthur/go-shortener/internal/config"
	"github.com/HungryArthur/go-shortener/internal/handler"
	"github.com/HungryArthur/go-shortener/internal/repository"
	"github.com/HungryArthur/go-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	config.Load()
	
	repo := repository.NewURLRepository()
	service := service.NewURLService(repo)
	handler := handler.NewURLHandler(service)

	router := chi.NewRouter()

	router.Get("/{shortenedURL}", handler.Get)
	router.Post("/", handler.Create)

	srv := http.Server{Addr: config.FlagRunAddr, Handler: router}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("can't start server", err)
		}
	}()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	quitCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := srv.Shutdown(quitCtx)
	if err != nil {
		fmt.Println("can't shutdown http server")
	} else {
		fmt.Println("successfully shutdowned http server")
	}

}
