package postgres

import (
	"context"
	"encoding/json"


	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/hihikaAAa/meeting-events/internal/app/ports"
)

type outboxRepo struct{ tx pgx.Tx }

func NewOutboxRepo(tx pgx.Tx) ports.OutboxRepository { 
  return &outboxRepo{tx: tx} 
}

func (r *outboxRepo) Add(ctx context.Context, aggregate string, aggregateID uuid.UUID, eventType string, payload any) error {
	b, err := json.Marshal(payload)
	if err != nil { 
    return err 
  }
	_, err = r.tx.Exec(ctx, `
		INSERT INTO outbox(aggregate, aggregate_id, event_type, payload, created_at)
		VALUES ($1,$2,$3,$4::jsonb, now())
	`, aggregate, aggregateID, eventType, b)
	return err
}


func (r *outboxRepo) FetchPending(ctx context.Context, limit int) ([]ports.OutboxEvent, error) {
	rows, err := r.tx.Query(ctx, `
		SELECT id, aggregate, aggregate_id, event_type, payload, created_at
		  FROM outbox
		 WHERE processed_at IS NULL
		 ORDER BY id
		 FOR UPDATE SKIP LOCKED
		 LIMIT $1
	`, limit)
	if err != nil { 
    return nil, err 
  }
	defer rows.Close()

	var out []ports.OutboxEvent
	for rows.Next() {
		var e ports.OutboxEvent
		if err := rows.Scan(&e.ID, &e.Aggregate, &e.AggregateID, &e.EventType, &e.Payload, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}


func (r *outboxRepo) MarkProcessed(ctx context.Context, ids []int64) error {
	if len(ids) == 0 { return nil }
	_, err := r.tx.Exec(ctx, `
		UPDATE outbox
		   SET processed_at = now()
		 WHERE id = ANY($1)
	`, ids)
	return err
}
