package meeting_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
	"log/slog"

	"github.com/docker/go-connections/nat"
	tc "github.com/testcontainers/testcontainers-go"
	tcp "github.com/testcontainers/testcontainers-go/modules/postgres"

	pg "github.com/hihikaAAa/meeting-events/internal/adapters/postgres"
	mig "github.com/hihikaAAa/meeting-events/internal/adapters/postgres/migrate"
	"github.com/hihikaAAa/meeting-events/internal/httpserver"
	hcreate "github.com/hihikaAAa/meeting-events/internal/httpserver/handlers/meeting/create"
	hdelete "github.com/hihikaAAa/meeting-events/internal/httpserver/handlers/meeting/delete"
	hget "github.com/hihikaAAa/meeting-events/internal/httpserver/handlers/meeting/get"
	hupdate "github.com/hihikaAAa/meeting-events/internal/httpserver/handlers/meeting/update"
	ucreate "github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/create"
	ucancel "github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/cancel"
	uupdate "github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/update"
	"github.com/hihikaAAa/meeting-events/internal/lib/logger/setup"
	
)

func TestMeetingsCRUD_E2E(t *testing.T) {
	t.Parallel()

	ctx := context.Background()


	pgC, err := tcp.RunContainer(ctx,
		tc.WithImage("postgres:16-alpine"),
		tcp.WithDatabase("meetings"),
		tcp.WithUsername("user"),
		tcp.WithPassword("pass"),
		tcp.WithInitScripts(), 
	)
	if err != nil {
		t.Fatalf("postgres up: %v", err)
	}
	t.Cleanup(func() { _ = pgC.Terminate(ctx) })

	host, _ := pgC.Host(ctx)
	mport, _ := pgC.MappedPort(ctx, nat.Port("5432/tcp"))
	dsn := "postgres://user:pass@" + host + ":" + mport.Port() + "/meetings?sslmode=disable"


	db, err := pg.New(ctx, dsn, 10, 30*time.Minute, 5*time.Minute)
	if err != nil {
		t.Fatalf("db init: %v", err)
	}
	t.Cleanup(func() { 
		db.Pool.Close() 
	})


	mdir := filepath.FromSlash("internal/adapters/postgres/init")
	if err := mig.Up(ctx, db.Pool, "file://"+mdir); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	log := setup.SetupLogger("test")
	slog.SetDefault(log)

	uow := pg.NewUoW(db.Pool)
	ucCreate := ucreate.New(uow)
	ucUpdate := uupdate.New(uow)
	ucCancel := ucancel.New(uow)

	handlers := httpserver.Handlers{
		Create: hcreate.New(log, ucCreate),
		Get:    hget.New(log, uow),
		Update: hupdate.New(log, ucUpdate),
		Delete: hdelete.New(log, ucCancel),
	}
	r := httpserver.NewRouter(handlers, log, "test", "test") 


	srv := httptest.NewServer(r)
	defer srv.Close()

	type createResp struct{ ID string `json:"id"` }

	body := `{"title":"Demo","starts_at":"2030-01-02T10:00:00Z","duration":90}`
	req, _ := http.NewRequest("POST", srv.URL+"/v1/meetings/", stringsReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("test", "test")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("create status=%d", res.StatusCode)
	}

	var cr createResp
	_ = json.NewDecoder(res.Body).Decode(&cr)
	if cr.ID == "" {
		t.Fatalf("empty id")
	}


	res2, err := http.Get(srv.URL + "/v1/meetings/" + cr.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res2.Body.Close()
	if res2.StatusCode != http.StatusOK {
		t.Fatalf("get status=%d", res2.StatusCode)
	}

	up := `{"title":"Updated","starts_at":"2030-01-02T12:00:00Z","duration":120}`
	req3, _ := http.NewRequest("PATCH", srv.URL+"/v1/meetings/"+cr.ID, stringsReader(up))
	req3.Header.Set("Content-Type", "application/json")
	req3.SetBasicAuth("test", "test")
	res3, err := http.DefaultClient.Do(req3)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	defer res3.Body.Close()
	if res3.StatusCode != http.StatusOK {
		t.Fatalf("update status=%d", res3.StatusCode)
	}

	req4, _ := http.NewRequest("DELETE", srv.URL+"/v1/meetings/"+cr.ID, nil)
	req4.SetBasicAuth("test", "test")
	res4, err := http.DefaultClient.Do(req4)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	defer res4.Body.Close()
	if res4.StatusCode != http.StatusOK {
		t.Fatalf("delete status=%d", res4.StatusCode)
	}
}

func stringsReader(s string) *os.File {
	tmp := filepath.Join(os.TempDir(), "b.json")
	_ = os.WriteFile(tmp, []byte(s), 0600)
	f, _ := os.Open(tmp)
	return f
}
