package delete

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log/slog"

	app "github.com/hihikaAAa/meeting-events/internal/app/app_errors"
	"github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/cancel"
	httpx "github.com/hihikaAAa/meeting-events/internal/httpserver/httpx"
)

type Handler struct {
	Log *slog.Logger
	UC  *cancel.UseCase
}

func New(log *slog.Logger, uc *cancel.UseCase) http.Handler {
	return &Handler{Log: log, UC: uc}
}

type response struct{ ID string `json:"id"` }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, app.ErrValidation.Error(), http.StatusBadRequest)
		return
	}

	out, err := h.UC.Handle(r.Context(), cancel.Input{ID: id})
	if err != nil {
		status, code, msg := httpx.HttpStatusFromErr(err)
    	httpx.WriteError(w, status, code, msg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response{ID: out.ID})
}
