package ports

import(
	"context"

	"github.com/google/uuid"

  	"github.com/hihikaAAa/meeting-events/internal/domain/meeting"
)

type MeetingRepository interface {
  Create(ctx context.Context, m *meeting.Meeting) error
  GetByID(ctx context.Context, id uuid.UUID) (*meeting.Meeting, error)
  Update(ctx context.Context, m *meeting.Meeting) error
  Cancel(ctx context.Context, id uuid.UUID) error
}

type OutboxRepository interface {
  Add(ctx context.Context, aggregate string, aggregateID uuid.UUID, eventType string, payload any) error
}