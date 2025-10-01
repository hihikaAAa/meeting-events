package postgres

import (
  "context"

  "github.com/google/uuid"
  "github.com/jackc/pgx/v5"
  "github.com/hihikaAAa/meeting-events/internal/app/ports"
)

type outboxRepo struct{ 
	tx pgx.Tx 
}

func NewOutboxRepo(tx pgx.Tx) ports.OutboxRepository {
	 return &outboxRepo{tx: tx}
	}

func (r *outboxRepo) Add(ctx context.Context, aggregate string, aggregateID uuid.UUID, eventType string, payload any) error {
  _, err := r.tx.Exec(ctx, `
    INSERT INTO outbox(aggregate, aggregate_id, event_type, payload)
    VALUES ($1,$2,$3,$4)
  `, aggregate, aggregateID, eventType, payload)
  return err
}
