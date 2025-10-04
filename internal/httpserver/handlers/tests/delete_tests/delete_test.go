package delete_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	tool "github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/tools"
	ucancel "github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/cancel"
	hdelete "github.com/hihikaAAa/meeting-events/internal/httpserver/handlers/meeting/delete"
	"github.com/hihikaAAa/meeting-events/internal/lib/logger/handlers/slogdiscard"
	"github.com/hihikaAAa/meeting-events/internal/domain/meeting"
)

func mountDelete(h http.Handler, id string) (*chi.Mux, string) {
	r := chi.NewRouter()
	r.Method(http.MethodDelete, "/v1/meetings/{id}", h)
	return r, "/v1/meetings/" + id
}

func TestDelete_HappyPath(t *testing.T) {
	log := slogdiscard.NewDiscardLogger()

	existing, _ := meeting.NewMeeting("Topic", time.Now().Add(time.Hour), 30*time.Minute)
	mr := &tool.MockMeetRepo{Fetched: existing}
	ob := &tool.MockOutbox{}
	uow := tool.FakeUoW{Repos: tool.FakeRepos{Mr: mr, Or: ob}}

	uc := ucancel.New(uow)
	h := hdelete.New(log, uc)

	r, url := mountDelete(h, existing.ID.String())
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.NotNil(t, mr.Updated)              
	require.Equal(t, 1, ob.Added)              
}

func TestDelete_BadID(t *testing.T) {
	log := slogdiscard.NewDiscardLogger()
	uc := ucancel.New(tool.FakeUoW{Repos: tool.FakeRepos{Mr: &tool.MockMeetRepo{}, Or: &tool.MockOutbox{}}})
	h := hdelete.New(log, uc)

	r, url := mountDelete(h, "not-a-uuid")
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDelete_Conflict(t *testing.T) {
	log := slogdiscard.NewDiscardLogger()

	existing, _ := meeting.NewMeeting("Topic", time.Now().Add(time.Hour), 30*time.Minute)
	_ = existing.Cancel() 
	mr := &tool.MockMeetRepo{Fetched: existing}
	uow := tool.FakeUoW{Repos: tool.FakeRepos{Mr: mr, Or: &tool.MockOutbox{}}}
	uc := ucancel.New(uow)
	h := hdelete.New(log, uc)

	r, url := mountDelete(h, existing.ID.String())
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusConflict, rr.Code)
}
