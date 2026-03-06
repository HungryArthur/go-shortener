package main

import (
	"net/http"

	"github.com/HungryArthur/go-shortener/internal/handler"
)

func main() {
	http.HandleFunc("/", handler.ShortPost)
	http.HandleFunc("/{id}", handler.ShortGet)
	

	http.ListenAndServe(":8080", nil)

}