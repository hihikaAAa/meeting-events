package integration

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/docker/go-connections/nat"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	pg "github.com/hihikaAAa/meeting-events/internal/adapters/postgres"
	mig "github.com/hihikaAAa/meeting-events/internal/adapters/postgres/migrate"
)

type PGTestEnv struct {
	DSN string
	DB  *pg.DB
	Ctx context.Context
}

func repoRoot() string {
	_, thisFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(thisFile), "../..")
}

func StartPostgres(tctx context.Context) (*PGTestEnv, func(), error) {
	pgC, err := postgres.Run(tctx,
		"postgres:16-alpine",
		postgres.WithDatabase("meetings"),
		postgres.WithUsername("user"),
		postgres.WithPassword("pass"),
		tc.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(2*time.Minute),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("postgres up: %w", err)
	}

	if _, err := pgC.MappedPort(tctx, nat.Port("5432/tcp")); err != nil {
		_ = pgC.Terminate(context.Background())
		return nil, nil, fmt.Errorf("port map: %w", err)
	}

	dsn, err := pgC.ConnectionString(tctx, "sslmode=disable")
	if err != nil {
		_ = pgC.Terminate(context.Background())
		return nil, nil, fmt.Errorf("dsn: %w", err)
	}
	dsn = strings.ReplaceAll(dsn, "@::1:", "@127.0.0.1:")
	dsn = strings.ReplaceAll(dsn, "@localhost:", "@127.0.0.1:")

	db, err := pg.New(tctx, dsn, 10, 30*time.Minute, 5*time.Minute)
	if err != nil {
		_ = pgC.Terminate(context.Background())
		return nil, nil, fmt.Errorf("db init: %w", err)
	}

	migDir := filepath.Join(repoRoot(), "internal", "adapters", "postgres", "init")
	if err := mig.Up(tctx, db.Pool, "file://"+migDir); err != nil {
		db.Pool.Close()
		_ = pgC.Terminate(context.Background())
		return nil, nil, fmt.Errorf("migrate: %w", err)
	}

	env := &PGTestEnv{
		DSN: dsn,
		DB:  db,
		Ctx: tctx,
	}
	teardown := func() {
		db.Pool.Close()
		_ = pgC.Terminate(context.Background())
	}
	return env, teardown, nil
}
