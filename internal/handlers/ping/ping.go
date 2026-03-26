package ping

import (
	"net/http"

	"github.com/HungryArthur/go-shortener/internal/service"
	"go.uber.org/zap"
)

type PingHandler struct {
	repo   service.URLRepository
	logger *zap.Logger
}

func NewPingHandler(repo service.URLRepository, logger *zap.Logger) *PingHandler {
	return &PingHandler{
		repo:   repo,
		logger: logger.With(zap.String("handler", "ping")),
	}
}

func (h *PingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.repo.Ping(); err != nil {
		h.logger.Error("database ping failed", zap.Error(err))
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
