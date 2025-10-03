package update

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log/slog"

	app "github.com/hihikaAAa/meeting-events/internal/app/app_errors"
	"github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/update"
	httpx "github.com/hihikaAAa/meeting-events/internal/httpserver/httpx"
)

type Handler struct {
	Log *slog.Logger
	UC  *update.UseCase
}

func New(log *slog.Logger, uc *update.UseCase) http.Handler {
	return &Handler{Log: log, UC: uc}
}

type request struct {
	Title    string        `json:"title"`    
	StartsAt time.Time     `json:"starts_at"`
	Duration int `json:"duration"` 
}

type response struct{ ID string `json:"id"` }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, app.ErrValidation.Error(), http.StatusBadRequest)
		return
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Log.Warn("bad json", slog.String("err", err.Error()))
    	httpx.WriteError(w, http.StatusBadRequest, "bad_json", "invalid request body")
		return
	}

	out, err := h.UC.Handle(r.Context(), update.Input{
		ID:       id,
		Title:    req.Title,
		StartsAt: req.StartsAt,
		Duration: time.Duration(req.Duration)*time.Minute,
	})
	if err != nil {
		status, code, msg := httpx.HttpStatusFromErr(err)
    	httpx.WriteError(w, status, code, msg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response{ID: out.ID})
}
