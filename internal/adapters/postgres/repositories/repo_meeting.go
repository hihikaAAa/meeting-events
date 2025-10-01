package postgres 

import (
  "context"
  "time"

  "github.com/google/uuid"
  "github.com/jackc/pgx/v5"
  app "github.com/hihikaAAa/meeting-events/internal/app/app_errors"
  "github.com/hihikaAAa/meeting-events/internal/app/ports"
  "github.com/hihikaAAa/meeting-events/internal/domain/meeting"
)

type meetingRepo struct{ 
	tx pgx.Tx 
}

func NewMeetingRepo(tx pgx.Tx) ports.MeetingRepository { 
	return &meetingRepo{tx: tx} 
}


func (r *meetingRepo) Create(ctx context.Context, m *meeting.Meeting) error {
  _, err := r.tx.Exec(ctx, `
    INSERT INTO meetings(id, title, starts_at, duration_sec, status, created_at, updated_at)
    VALUES ($1,$2,$3,$4,$5, now(), now())
  `, m.ID, m.Title, m.StartsAt, int(m.Duration/time.Second), m.Status)
  return err
}

func (r *meetingRepo) GetByID(ctx context.Context, id uuid.UUID) (*meeting.Meeting, error) {
  row := r.tx.QueryRow(ctx, `
    SELECT id, title, starts_at, duration_sec, status, created_at, updated_at
    FROM meetings WHERE id = $1
  `, id)

  var (
    mm meeting.Meeting 
	durSec int
    createdAt, updatedAt time.Time
  )

  if err := row.Scan(&mm.ID, &mm.Title, &mm.StartsAt, &durSec, &mm.Status, &createdAt, &updatedAt); err != nil {
    if err == pgx.ErrNoRows { 
		return nil, app.ErrNotFound
	}
    return nil, err
  }
  mm.Duration = time.Duration(durSec) * time.Second
  mm.RestoreTimestamps(createdAt, updatedAt)
  return &mm, nil
}

func (r *meetingRepo) Update(ctx context.Context, m *meeting.Meeting) error {
  _, err := r.tx.Exec(ctx, `
    UPDATE meetings
       SET title=$2, starts_at=$3, duration_sec=$4, status=$5, updated_at=now()
     WHERE id=$1
  `, m.ID, m.Title, m.StartsAt, int(m.Duration/time.Second), m.Status)
  return err
}

func (r *meetingRepo) Cancel(ctx context.Context, id uuid.UUID) error {
  cmd, err := r.tx.Exec(ctx, `UPDATE meetings SET status='canceled', updated_at=now() WHERE id=$1`, id)
  if err != nil { return err }
  if cmd.RowsAffected() == 0 { return app.ErrNotFound }
  return nil
}