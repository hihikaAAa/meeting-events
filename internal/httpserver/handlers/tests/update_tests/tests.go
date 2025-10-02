package update_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	tool "github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/tools"
	uupdate "github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/update"
	hupdate "github.com/hihikaAAa/meeting-events/internal/httpserver/handlers/meeting/update"
	"github.com/hihikaAAa/meeting-events/internal/lib/logger/handlers/slogdiscard"
	"github.com/hihikaAAa/meeting-events/internal/domain/meeting"
)

func mountUpdate(h http.Handler, id string) (*chi.Mux, string) {
	r := chi.NewRouter()
	r.Method(http.MethodPatch, "/v1/meetings/{id}", h)
	return r, "/v1/meetings/" + id
}

func TestUpdate_HappyPath(t *testing.T) {
	log := slogdiscard.NewDiscardLogger()

	existing, _ := meeting.NewMeeting("Old", time.Now().Add(time.Hour), 30*time.Minute)
	mr := &tool.MockMeetRepo{Fetched: existing}
	ob := &tool.MockOutbox{}
	uow := tool.FakeUoW{Repos: tool.FakeRepos{Mr: mr, Or: ob}}

	uc := uupdate.New(uow)
	h := hupdate.New(log, uc)

	body := map[string]any{
		"title":    "New",
		"duration": (90 * time.Minute).Nanoseconds(),
	}
	b, _ := json.Marshal(body)

	r, url := mountUpdate(h, existing.ID.String())
	req := httptest.NewRequest(http.MethodPatch, url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.NotNil(t, mr.Updated)
	require.Equal(t, "New", mr.Updated.Title)
	require.Equal(t, 1, ob.Added)
}

func TestUpdate_BadID(t *testing.T) {
	log := slogdiscard.NewDiscardLogger()
	uc := uupdate.New(tool.FakeUoW{Repos: tool.FakeRepos{Mr: &tool.MockMeetRepo{}, Or: &tool.MockOutbox{}}})
	h := hupdate.New(log, uc)

	r, url := mountUpdate(h, "not-a-uuid")
	req := httptest.NewRequest(http.MethodPatch, url, bytes.NewReader([]byte(`{}`)))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdate_Validation(t *testing.T) {
	log := slogdiscard.NewDiscardLogger()

	existing, _ := meeting.NewMeeting("Old", time.Now().Add(time.Hour), 30*time.Minute)
	mr := &tool.MockMeetRepo{Fetched: existing}
	uow := tool.FakeUoW{Repos: tool.FakeRepos{Mr: mr, Or: &tool.MockOutbox{}}}
	uc := uupdate.New(uow)
	h := hupdate.New(log, uc)

	body := map[string]any{"title": ""} 
	b, _ := json.Marshal(body)

	r, url := mountUpdate(h, existing.ID.String())
	req := httptest.NewRequest(http.MethodPatch, url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}
