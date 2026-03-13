package handler

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/HungryArthur/go-shortener/internal/service"
)

type URLService interface {
	Save(url string) (string, error)
	Get(string) (string, error)
}

type URLHandler struct {
	service URLService
}

func NewURLHandler(service URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

func (h *URLHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *URLHandler) Get(w http.ResponseWriter, r *http.Request) {
	shortenedURL, _ := strings.CutPrefix(r.URL.Path, "/")
	srcURL, err := h.service.Get(shortenedURL)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrShortenedURLDoesntExist):
			http.NotFound(w, r)
		case errors.Is(err, service.ErrCantGetURL):
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Location", srcURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *URLHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	body, err := io.ReadAll(io.LimitReader(r.Body, 100))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("can't read request's body"))
		return
	}
	if n, _ := r.Body.Read(make([]byte, 1)); n == 1 {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		w.Write([]byte("request's body too big"))
		return
	}

	urlIn := string(body)

	shortenedURL, err := h.service.Save(urlIn)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can't create url"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + shortenedURL))
}
