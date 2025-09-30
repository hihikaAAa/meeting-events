package create_test

import(
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	app "github.com/hihikaAAa/meeting-events/internal/app/app_errors"
	"github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/create"
	tool "github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/tools"
)

func TestCreate_HappyPath(t *testing.T) {
	mr := &tool.MockMeetRepo{}
	ob := &tool.MockOutbox{}
	uow := tool.FakeUoW{Repos: tool.FakeRepos{Mr: mr, Or: ob}}

	uc := create.New(uow)
	out, err := uc.Handle(context.Background(), create.Input{
		Title: "Demo",
		StartsAt: time.Now().Add(time.Hour),
		Duration: 45 * time.Minute,
	})
	require.NoError(t, err)
	require.NotEmpty(t, out.ID)
	require.NotNil(t, mr.Created)
	require.Equal(t, 1, ob.Added) 
}

func TestCreate_Validation(t *testing.T) {
	mr := &tool.MockMeetRepo{}
	ob := &tool.MockOutbox{}
	uow := tool.FakeUoW{Repos: tool.FakeRepos{Mr: mr, Or: ob}}

	uc := create.New(uow)
	_, err := uc.Handle(context.Background(), create.Input{
		Title: "",
		StartsAt: time.Now().Add(time.Hour),
		Duration: 30*time.Minute,
	})
	require.ErrorIs(t, err, app.ErrValidation)
}

