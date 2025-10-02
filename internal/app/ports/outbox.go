package ports

import(
	"time"
	"github.com/google/uuid"
)

type OutboxEvent struct {
	ID          int64
	Aggregate   string
	AggregateID uuid.UUID
	EventType   string
	Payload     []byte
	CreatedAt   time.Time
}