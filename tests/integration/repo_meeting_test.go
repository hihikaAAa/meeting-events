package integration

import (
	"context"
	"testing"
	"time"

	rep "github.com/hihikaAAa/meeting-events/internal/adapters/postgres/repositories"
	"github.com/hihikaAAa/meeting-events/internal/domain/meeting"
)

func TestMeetingRepo_CRUD(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	env, done, err := StartPostgres(ctx)
	if err != nil { t.Fatal(err) }
	defer done()

	tx, err := env.DB.Pool.Begin(ctx)
	if err != nil { 
		t.Fatal(err) 
	}
	defer tx.Rollback(ctx) 

	repo := rep.NewMeetingRepo(tx)

	m, err := meeting.NewMeeting("Title", time.Now().Add(1*time.Hour), 90*time.Minute)
	if err != nil {
		 t.Fatal(err)
		}
	if err := repo.Create(ctx, m); err != nil {
		 t.Fatal(err) 
		}

	got, err := repo.GetByID(ctx, m.ID)
	if err != nil { 
		t.Fatal(err) 
	}
	if got.Title != "Title" { 
		t.Fatalf("title mismatch: %q", got.Title) 
	}

	if err := m.Update("Upd", time.Now().Add(2*time.Hour), 120*time.Minute); err != nil {
		 t.Fatal(err) 
		}
	if err := repo.Update(ctx, m); err != nil {
		 t.Fatal(err) 
		}
	got2, err := repo.GetByID(ctx, m.ID)
	if err != nil { 
		t.Fatal(err) 
	}
	if got2.Title != "Upd" { 
		t.Fatalf("title mismatch after update: %q", got2.Title)
	 }

	if err := repo.Cancel(ctx, m.ID); err != nil { 
		t.Fatal(err) 
	}

	var status string
	if err := tx.QueryRow(ctx, `select status from meetings where id=$1`, m.ID).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != string(meeting.StatusCanceled) {
		t.Fatalf("expected canceled, got %s", status)
	}
}
