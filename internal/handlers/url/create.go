package url

import (
	"io"
	"net/http"
)

func (h *URLHandler) CreateTextPlain(w http.ResponseWriter, r *http.Request) {
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
	w.Write([]byte(shortenedURL))
}
