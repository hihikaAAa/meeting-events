package gettests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/hihikaAAa/meeting-events/internal/app/app_errors"
	tool "github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/tools"
	hget "github.com/hihikaAAa/meeting-events/internal/httpserver/handlers/meeting/get"
	"github.com/hihikaAAa/meeting-events/internal/lib/logger/handlers/slogdiscard"
	"github.com/hihikaAAa/meeting-events/internal/domain/meeting"
)

func mountGet(h http.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/v1/meetings/{id}", h)
	return r
}

func TestGet_HappyPath(t *testing.T) {
	log := slogdiscard.NewDiscardLogger()

	m, _ := meeting.NewMeeting("Topic", time.Now().Add(time.Hour), 30*time.Minute)
	mr := &tool.MockMeetRepo{Fetched: m}
	uow := tool.FakeUoW{Repos: tool.FakeRepos{Mr: mr, Or: &tool.MockOutbox{}}}

	h := hget.New(log, uow)

	req := httptest.NewRequest(http.MethodGet, "/v1/meetings/"+m.ID.String(), nil)
	rr := httptest.NewRecorder()

	mountGet(h).ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Contains(t, rr.Body.String(), m.ID.String())
}

func TestGet_BadID(t *testing.T) {
	log := slogdiscard.NewDiscardLogger()
	uow := tool.FakeUoW{Repos: tool.FakeRepos{Mr: &tool.MockMeetRepo{}, Or: &tool.MockOutbox{}}}
	h := hget.New(log, uow)

	req := httptest.NewRequest(http.MethodGet, "/v1/meetings/not-a-uuid", nil)
	rr := httptest.NewRecorder()

	mountGet(h).ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGet_NotFound(t *testing.T) {
	log := slogdiscard.NewDiscardLogger()
	mr := &tool.MockMeetRepo{GetErr: app.ErrNotFound} 
	uow := tool.FakeUoW{Repos: tool.FakeRepos{Mr: mr, Or: &tool.MockOutbox{}}}
	h := hget.New(log, uow)

	req := httptest.NewRequest(http.MethodGet, "/v1/meetings/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()

	mountGet(h).ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}
