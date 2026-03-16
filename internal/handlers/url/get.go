package url

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/HungryArthur/go-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

func (h *URLHandler) GetTextPlain(w http.ResponseWriter, r *http.Request) {
	shortenedURL := chi.URLParam(r, "shortenedURL")
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

func (h *URLHandler) GetJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	dtoIn := GetRequest{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&dtoIn); err != nil {
		jsonErrResp(w, http.StatusBadRequest, "invalid json")
		return
	}
	srcURL, err := h.service.Get(dtoIn.ShortenedURL)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrShortenedURLDoesntExist):
			jsonErrResp(w, http.StatusNotFound, "shortened url not foudn")
		case errors.Is(err, service.ErrCantGetURL):
			jsonErrResp(w, http.StatusInternalServerError, "can't get url")
		}
		return
	}
	dtoOut := GetResponse{
		SourceURL: srcURL,
	}
	err = json.NewEncoder(w).Encode(dtoOut)
	if err != nil {
		// logging
		jsonErrResp(w, http.StatusInternalServerError, "can't get url")
		return
	}
}
