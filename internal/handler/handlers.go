package handler

import (
	// "fmt"
	"net/http"
)

func ShortPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST request allowed", http.StatusMethodNotAllowed)
		return
	}

	// body := fmt.Sprintf("Method: %s\r\n", r.Method)


	// r.Write([]byte(body))

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")

}
