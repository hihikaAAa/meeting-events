package update_test

import(
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	app "github.com/hihikaAAa/meeting-events/internal/app/app_errors"
	"github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/update"
	tool "github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/tools"
	"github.com/hihikaAAa/meeting-events/internal/domain/meeting"
)

func TestUpdate_HappyPath(t *testing.T) {
	existing, err := meeting.NewMeeting("Old", time.Now().Add(time.Hour), 30*time.Minute)
	require.NoError(t, err)

	mr := &tool.MockMeetRepo{Fetched: existing}
	ob := &tool.MockOutbox{}
	uow := tool.FakeUoW{Repos: tool.FakeRepos{Mr: mr, Or: ob}}

	uc := update.New(uow)

	newTitle := "New"
	newDuration := 90 * time.Minute

	out, err := uc.Handle(context.Background(), update.Input{
		ID: existing.ID.String(),
		Title: newTitle,        
		Duration: newDuration,   
	})
	require.NoError(t, err)
	require.Equal(t, existing.ID.String(), out.ID)

	require.NotNil(t, mr.Updated)
	require.Equal(t, "New", mr.Updated.Title)
	require.Equal(t, newDuration, mr.Updated.Duration)

	require.Equal(t, 1, ob.Added)
}

func TestUpdate_NotFound(t *testing.T) {
	mr := &tool.MockMeetRepo{GetErr: app.ErrNotFound} 
	ob := &tool.MockOutbox{}
	uow := tool.FakeUoW{Repos: tool.FakeRepos{Mr: mr, Or: ob}}

	uc := update.New(uow)

	newTitle := "Whatever"
	_, err := uc.Handle(context.Background(), update.Input{
		ID: "8c0f5f7d-8b1e-4b2d-9a3e-0a9e8b6d1a11",
		Title: newTitle,
	})
	require.ErrorIs(t, err, app.ErrNotFound)
	require.Nil(t, mr.Updated)
	require.Equal(t, 0, ob.Added)
}
