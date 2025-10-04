package outbox_test

import(
	"time"
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hihikaAAa/meeting-events/internal/app/ports"
	obw "github.com/hihikaAAa/meeting-events/internal/services/outboxworker"
	"github.com/hihikaAAa/meeting-events/internal/lib/logger/handlers/slogdiscard"
	tool "github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/tools"
)

func TestWorker_Tick_MarksProcessed(t *testing.T) {
	log := slogdiscard.NewDiscardLogger()

	pending := []ports.OutboxEvent{{ID: 1, EventType: "MeetingCreated"}, {ID: 2, EventType: "MeetingUpdated"}}
	mr := &tool.MockMeetRepo{}
	ob := &tool.MockOutbox{Pending: pending}
	uow := tool.FakeUoW{Repos: tool.FakeRepos{Mr: mr, Or: ob}}

	pub := obw.NewLogPublisher(log)
	w := obw.New(log, uow, pub, 10, time.Millisecond)

	require.NoError(t, w.Tick(context.Background()))
	require.Equal(t, []int64{1, 2}, ob.Marked)
}

