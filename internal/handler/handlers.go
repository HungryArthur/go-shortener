package handler

import (
	// "fmt"
	"io"
	"net/http"
)

type UrlService interface {
	Save(url string) (string, error)
}

type UrlHandler struct {
	service UrlService
}

func NewUrlHandler(service UrlService) *UrlHandler {
	return &UrlHandler{
		service: service,
	}
}

func (h *UrlHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST request allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("can't read request's body"))
		return
	}

	urlIn := string(body)

	shortenedUrl, err := h.service.Save(urlIn)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can't create url"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortenedUrl))
}
