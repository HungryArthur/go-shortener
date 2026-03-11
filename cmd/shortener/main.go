package main

import (
	"net/http"

	"github.com/HungryArthur/go-shortener/internal/handler"
	"github.com/HungryArthur/go-shortener/internal/repository"
	"github.com/HungryArthur/go-shortener/internal/service"
)

func main() {
	repo := repository.NewUrlRepository()
	service := service.NewUrlService(repo)
	handler := handler.NewUrlHandler(service)

	http.HandleFunc("/", handler.Handle)

	http.ListenAndServe(":8080", nil)
}
