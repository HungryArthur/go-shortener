package main

import (
	"fmt"
	"net/http"

	"github.com/HungryArthur/go-shortener/internal/handler"
	"github.com/HungryArthur/go-shortener/internal/repository"
	"github.com/HungryArthur/go-shortener/internal/service"
)

func main() {
	repo := repository.NewURLRepository()
	service := service.NewURLService(repo)
	handler := handler.NewURLHandler(service)

	http.HandleFunc("/", handler.Handle)

	err := http.ListenAndServe(":8080", nil)
	if err != nil && err != http.ErrServerClosed {
		fmt.Println("can't start server", err)
	}
}
