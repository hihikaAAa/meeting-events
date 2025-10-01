package postgres

import (
  "context"

  "github.com/jackc/pgx/v5"
  "github.com/jackc/pgx/v5/pgxpool"
  "github.com/hihikaAAa/meeting-events/internal/app/ports"
  rep "github.com/hihikaAAa/meeting-events/internal/adapters/postgres/repositories"
)

type uow struct {
  pool *pgxpool.Pool
}

func NewUoW(pool *pgxpool.Pool) ports.UnitOfWork {
  return &uow{pool: pool}
}

type repos struct {
  tx pgx.Tx
}

func (r repos) Meetings() ports.MeetingRepository { 
	return rep.NewMeetingRepo(r.tx) 
}
func (r repos) Outbox() ports.OutboxRepository   {
	 return rep.NewOutboxRepo(r.tx)
}

func (u *uow) WithinTx(ctx context.Context, fn func(ctx context.Context, r ports.Repos) error) error {
  tx, err := u.pool.BeginTx(ctx, pgx.TxOptions{})
  if err != nil { 
	return err
	}
  r := repos{tx: tx}
  if err := fn(ctx, r); err != nil {
    _ = tx.Rollback(ctx)
    return err
  }
  return tx.Commit(ctx)
}
