package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"runtime"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	pg "github.com/hihikaAAa/meeting-events/internal/adapters/postgres"
	mig "github.com/hihikaAAa/meeting-events/internal/adapters/postgres/migrate"

	ucancel "github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/cancel"
	ucreate "github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/create"
	uupdate "github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/update"

	"log/slog"

	"github.com/hihikaAAa/meeting-events/internal/httpserver"
	hcreate "github.com/hihikaAAa/meeting-events/internal/httpserver/handlers/meeting/create"
	hdelete "github.com/hihikaAAa/meeting-events/internal/httpserver/handlers/meeting/delete"
	hget "github.com/hihikaAAa/meeting-events/internal/httpserver/handlers/meeting/get"
	hupdate "github.com/hihikaAAa/meeting-events/internal/httpserver/handlers/meeting/update"
	"github.com/hihikaAAa/meeting-events/internal/lib/logger/setup"
)

func TestMeetingsCRUD_E2E(t *testing.T) {


	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	 pgC, err := postgres.Run(ctx,
        "postgres:16-alpine",
        postgres.WithDatabase("meetings"),
        postgres.WithUsername("user"),
        postgres.WithPassword("pass"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).
                WithStartupTimeout(2 * time.Minute),
        ),
    )
	if err != nil {
		t.Fatalf("postgres up: %v", err)
	}
	t.Cleanup(func() { _ = pgC.Terminate(context.Background()) })

	if _, err := pgC.MappedPort(ctx, nat.Port("5432/tcp")); err != nil {
		t.Fatalf("port map: %v", err)
	}

	dsn, err := pgC.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("dsn: %v", err)																		
	}

	dsn = strings.ReplaceAll(dsn, "@::1:", "@127.0.0.1:")
	dsn = strings.ReplaceAll(dsn, "@localhost:", "@127.0.0.1:")

	db, err := pg.New(ctx, dsn, 10, 30*time.Minute, 5*time.Minute)
	if err != nil {
		t.Fatalf("db init: %v", err)
	}
	t.Cleanup(func() { db.Pool.Close() })

	migDir := filepath.Join(repoRoot(), "internal", "adapters", "postgres", "init")
	if err := mig.Up(ctx, db.Pool, "file://" + migDir); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	log := setup.SetupLogger("test")
	slog.SetDefault(log)

	uow := pg.NewUoW(db.Pool)
	ucC := ucreate.New(uow)
	ucU := uupdate.New(uow)
	ucD := ucancel.New(uow)

	handlers := httpserver.Handlers{
		Create: hcreate.New(log, ucC),
		Get:    hget.New(log, uow),
		Update: hupdate.New(log, ucU),
		Delete: hdelete.New(log, ucD),
	}
	r := httpserver.NewRouter(handlers, log, "test", "test")

	srv := httptest.NewServer(r)
	defer srv.Close()

	type createResp struct{ ID string `json:"id"` }

	reqBody := `{"title":"Demo","starts_at":"2030-01-02T10:00:00Z","duration":90}`
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/v1/meetings/", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("test", "test")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("create: %v", err)
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

	getRes, err := http.Get(srv.URL + "/v1/meetings/" + cr.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("get status=%d", getRes.StatusCode)
	}

	upBody := `{"title":"Updated","starts_at":"2030-01-02T12:00:00Z","duration":120}`
	upReq, _ := http.NewRequest(http.MethodPatch, srv.URL+"/v1/meetings/"+cr.ID, strings.NewReader(upBody))
	upReq.Header.Set("Content-Type", "application/json")
	upReq.SetBasicAuth("test", "test")

	upRes, err := http.DefaultClient.Do(upReq)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	defer upRes.Body.Close()
	if upRes.StatusCode != http.StatusOK {
		t.Fatalf("update status=%d", upRes.StatusCode)
	}

	delReq, _ := http.NewRequest(http.MethodDelete, srv.URL+"/v1/meetings/"+cr.ID, nil)
	delReq.SetBasicAuth("test", "test")
	delRes, err := http.DefaultClient.Do(delReq)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	defer delRes.Body.Close()
	if delRes.StatusCode != http.StatusOK {
		t.Fatalf("delete status=%d", delRes.StatusCode)
	}
}


func repoRoot() string {
    _, thisFile, _, _ := runtime.Caller(0)
    return filepath.Join(filepath.Dir(thisFile), "..")
}