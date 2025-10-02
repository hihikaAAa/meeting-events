package create

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/create"
	httpx "github.com/hihikaAAa/meeting-events/internal/httpserver/httpx"
)

type Handler struct {
	Log *slog.Logger
	UC  *create.UseCase
}

type request struct {
	Title    string        `json:"title"`
	StartsAt time.Time     `json:"starts_at"`
	Duration time.Duration `json:"duration"` 
}

type response struct {
	ID string `json:"id"`
}

func New(log *slog.Logger, uc *create.UseCase) http.Handler {
	return &Handler{Log: log, UC: uc}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest); return
	}
	out, err := h.UC.Handle(r.Context(), create.Input{
		Title:    req.Title,
		StartsAt: req.StartsAt,
		Duration: req.Duration,
	})
	if err != nil {
		http.Error(w, err.Error(), httpx.HttpStatusFromErr(err)); return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response{ID: out.ID})
}


