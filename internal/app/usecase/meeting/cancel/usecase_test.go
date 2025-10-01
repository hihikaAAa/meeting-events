package cancel_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	app "github.com/hihikaAAa/meeting-events/internal/app/app_errors"
	"github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/cancel"
	tool "github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/tools"
	"github.com/hihikaAAa/meeting-events/internal/domain/meeting"
)

func TestCancel_HappyPath(t *testing.T) {
	existing, err := meeting.NewMeeting("Demo", time.Now().Add(time.Hour), 45*time.Minute)
	require.NoError(t, err)

	mr := &tool.MockMeetRepo{Fetched: existing}
	ob := &tool.MockOutbox{}
	uow := tool.FakeUoW{Repos: tool.FakeRepos{Mr: mr, Or: ob}}

	uc := cancel.New(uow)

	out, err := uc.Handle(context.Background(), cancel.Input{
		ID: existing.ID.String(),
	})
	require.NoError(t, err)
	require.Equal(t, existing.ID.String(), out.ID)

	require.NotNil(t, mr.Updated)
	require.Equal(t, meeting.StatusCanceled, mr.Updated.Status)
	require.Equal(t, 1, ob.Added)
}

func TestCancel_AlreadyCanceled(t *testing.T) {
	existing, err := meeting.NewMeeting("Demo", time.Now().Add(time.Hour), 30*time.Minute)
	require.NoError(t, err)
	require.NoError(t, existing.Cancel())

	mr := &tool.MockMeetRepo{Fetched: existing}
	ob := &tool.MockOutbox{}
	uow := tool.FakeUoW{Repos: tool.FakeRepos{Mr: mr, Or: ob}}

	uc := cancel.New(uow)

	_, err = uc.Handle(context.Background(), cancel.Input{
		ID: existing.ID.String(),
	})
	require.ErrorIs(t, err, app.ErrConflict)
	require.Nil(t, mr.Updated)
	require.Equal(t, 0, ob.Added)
}

func TestCancel_Ongoing(t *testing.T) {
	existing, err := meeting.NewMeeting("Demo", time.Now().Add(time.Hour), 30*time.Minute)
	require.NoError(t, err)
	existing.Status = meeting.StatusOngoing

	mr := &tool.MockMeetRepo{Fetched: existing}
	ob := &tool.MockOutbox{}
	uow := tool.FakeUoW{Repos: tool.FakeRepos{Mr: mr, Or: ob}}

	uc := cancel.New(uow)

	_, err = uc.Handle(context.Background(), cancel.Input{
		ID: existing.ID.String(),
	})
	require.ErrorIs(t, err, app.ErrConflict)
	require.Nil(t, mr.Updated)
	require.Equal(t, 0, ob.Added)
}

func TestCancel_NotFound(t *testing.T) {
	mr := &tool.MockMeetRepo{GetErr: app.ErrNotFound}
	ob := &tool.MockOutbox{}
	uow := tool.FakeUoW{Repos: tool.FakeRepos{Mr: mr, Or: ob}}

	uc := cancel.New(uow)

	_, err := uc.Handle(context.Background(), cancel.Input{
		ID: "8c0f5f7d-8b1e-4b2d-9a3e-0a9e8b6d1a11",
	})
	require.ErrorIs(t, err, app.ErrNotFound)
	require.Nil(t, mr.Updated)
	require.Equal(t, 0, ob.Added)
}

func TestCancel_BadID(t *testing.T) {
	mr := &tool.MockMeetRepo{}
	ob := &tool.MockOutbox{}
	uow := tool.FakeUoW{Repos: tool.FakeRepos{Mr: mr, Or: ob}}

	uc := cancel.New(uow)

	_, err := uc.Handle(context.Background(), cancel.Input{
		ID: "not-a-uuid",
	})
	require.ErrorIs(t, err, app.ErrValidation)
	require.Nil(t, mr.Updated)
	require.Equal(t, 0, ob.Added)
}
