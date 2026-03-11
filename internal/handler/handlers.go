package handler

import (
	// "fmt"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/HungryArthur/go-shortener/internal/service"
)

type UrlService interface {
	Save(url string) (string, error)
	Get(string) (string, error)
}

type UrlHandler struct {
	service UrlService
}

func NewUrlHandler(service UrlService) *UrlHandler {
	return &UrlHandler{
		service: service,
	}
}

func (h *UrlHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *UrlHandler) Get(w http.ResponseWriter, r *http.Request) {
	shortenedUrl, _ := strings.CutPrefix(r.URL.Path, "/")
	srcUrl, err := h.service.Get(shortenedUrl)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrShortenedUrlDoesntExist):
			http.NotFound(w, r)
		case errors.Is(err, service.ErrCantGetUrl):
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Location", srcUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
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
