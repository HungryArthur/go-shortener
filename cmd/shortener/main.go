package main

import (
	"fmt"
	"net/http"

	"github.com/HungryArthur/go-shortener/internal/handler"
	"github.com/HungryArthur/go-shortener/internal/repository"
	"github.com/HungryArthur/go-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	repo := repository.NewURLRepository()
	service := service.NewURLService(repo)
	handler := handler.NewURLHandler(service)

	router := chi.NewRouter()

	router.Get("/{shortenedPath}", handler.Get)
	router.Post("/", handler.Create)

	srv := http.Server{Addr: ":8080", Handler: router}

	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Println("can't start server", err)
	}
}
