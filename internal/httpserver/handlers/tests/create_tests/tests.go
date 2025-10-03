package createtests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/create"
	tool "github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/tools"
	hcreate "github.com/hihikaAAa/meeting-events/internal/httpserver/handlers/meeting/create"
	"github.com/hihikaAAa/meeting-events/internal/lib/logger/handlers/slogdiscard"
)

func router(h http.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Method(http.MethodPost, "/", h)
	return r
}

func TestCreate_HappyPath(t *testing.T) {
	log := slogdiscard.NewDiscardLogger()
	mr := &tool.MockMeetRepo{}
	ob := &tool.MockOutbox{}
	uow := tool.FakeUoW{Repos: tool.FakeRepos{Mr: mr, Or: ob}}
	uc := create.New(uow)

	h := hcreate.New(log, uc)

	body := map[string]any{
		"title":     "Demo",
		"starts_at": time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
		"duration":  (45 * time.Minute),
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router(h).ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.NotEmpty(t, rr.Body.String())
	require.NotNil(t, mr.Created)
	require.Equal(t, 1, ob.Added)
}

func TestCreate_BadJSON(t *testing.T) {
	log := slogdiscard.NewDiscardLogger()
	uow := tool.FakeUoW{Repos: tool.FakeRepos{Mr: &tool.MockMeetRepo{}, Or: &tool.MockOutbox{}}}
	uc := create.New(uow)
	h := hcreate.New(log, uc)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{bad"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router(h).ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreate_ValidationError(t *testing.T) {
	log := slogdiscard.NewDiscardLogger()
	uow := tool.FakeUoW{Repos: tool.FakeRepos{Mr: &tool.MockMeetRepo{}, Or: &tool.MockOutbox{}}}
	uc := create.New(uow)
	h := hcreate.New(log, uc)

	body := map[string]any{
		"title":     "", 						
		"starts_at": time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
		"duration":  (30 * time.Minute),
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router(h).ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}
