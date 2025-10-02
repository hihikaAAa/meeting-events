package get

import (
	"encoding/json"
	"net/http"
	"time"
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log/slog"

	app "github.com/hihikaAAa/meeting-events/internal/app/app_errors"
	"github.com/hihikaAAa/meeting-events/internal/app/ports"
	httpx "github.com/hihikaAAa/meeting-events/internal/httpserver/httpx"
)

type Handler struct {
	Log *slog.Logger
	UoW ports.UnitOfWork
}

func New(log *slog.Logger, uow ports.UnitOfWork) http.Handler {
	return &Handler{Log: log, UoW: uow}
}

type response struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	StartsAt  string `json:"starts_at"`
	DurationS int64  `json:"duration_seconds"`
	Status    string `json:"status"`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, app.ErrValidation.Error(), http.StatusBadRequest)
		return
	}

	var resp response
	err = h.UoW.WithinTx(r.Context(), func(ctx context.Context, repos ports.Repos) error {
		m, err := repos.Meetings().GetByID(ctx, id)
		if err != nil {
			return app.ErrNotFound
		}
		resp = response{
			ID:        m.ID.String(),
			Title:     m.Title,
			StartsAt:  m.StartsAt.UTC().Format(time.RFC3339),
			DurationS: int64(m.Duration / time.Second),
			Status:    string(m.Status),
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), httpx.HttpStatusFromErr(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
