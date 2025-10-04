package integration

import (
	"context"
	"testing"
	"time"
	"fmt"

	"github.com/google/uuid"
	pguow "github.com/hihikaAAa/meeting-events/internal/adapters/postgres"
	"github.com/hihikaAAa/meeting-events/internal/app/ports"
	"github.com/hihikaAAa/meeting-events/internal/domain/meeting"
)


func TestUoW_CommitRollback(t *testing.T) {
	var assertErr = fmt.Errorf("rollback")
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	env, done, err := StartPostgres(ctx)
	if err != nil {
		 t.Fatal(err) 
		}
	defer done()

	uow := pguow.NewUoW(env.DB.Pool)

	var createdID uuid.UUID
	err = uow.WithinTx(ctx, func(ctx context.Context, r ports.Repos) error {
		m, _ := meeting.NewMeeting("C", time.Now().Add(1*time.Hour), 60*time.Minute)
		createdID = m.ID
		return r.Meetings().Create(ctx, m)
	})
	if err != nil {
		 t.Fatalf("commit case: %v", err) 
		}

	err = uow.WithinTx(ctx, func(ctx context.Context, r ports.Repos) error {
		_, err := r.Meetings().GetByID(ctx, createdID)
		return err
	})
	if err != nil { t.Fatalf("must exist after commit: %v", err) }

	var toRollback uuid.UUID
	_ = uow.WithinTx(ctx, func(ctx context.Context, r ports.Repos) error {
		m, _ := meeting.NewMeeting("R", time.Now().Add(1*time.Hour), 60*time.Minute)
		toRollback = m.ID
		_ = r.Meetings().Create(ctx, m)
		return assertErr
	})

	err = uow.WithinTx(ctx, func(ctx context.Context, r ports.Repos) error {
		_, err := r.Meetings().GetByID(ctx, toRollback)
		if err == nil {
			return fmt.Errorf("expected not found after rollback")
		}
		return nil
	})
	if err != nil { t.Fatal(err) }
}

