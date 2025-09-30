package config_test 

import(
	"os"
	"path/filepath"
	"testing"
	"time"
	"os/exec"

	"github.com/hihikaAAa/meeting-events/internal/config"
)

func writeTempConfig(t * testing.T, body string) string{
	t.Helper()
	dir:= t.TempDir()
	p := filepath.Join(dir,"cfg.yaml")
	if err := os.WriteFile(p,[]byte(body),0o600); err != nil{
		t.Fatalf("write temp cfg: %v", err)
	}
	return p
}

func TestMustLoad_ReadsYAML(t *testing.T) {
	yaml := `
env: "local"
app:
  name: "meeting-svc"
  http:
    address: "localhost:8081"
    timeouts:
      read: 4s
      write: 6s
      idle: 60s
      event: 5s
    user: "hihika"
    password: "pwd"
db:
  dsn: "postgres://u:p@db:5432/meetings?sslmode=disable"
  max_open_conns: 20
  max_idle_conns: 5
  conn_max_lifetime: "30m"
migrations: { dir: "file://migrations" }
outbox: { poll_interval: "3s", batch_size: 100 }
`
	p := writeTempConfig(t, yaml)
	t.Setenv("CONFIG_PATH", p)

	cfg := config.MustLoad()

	if cfg.Env != "local" {
		t.Fatalf("env want local, got %s", cfg.Env)
	}
	if cfg.App.HTTP.Address != "localhost:8081" {
		t.Fatalf("address mismatch: %s", cfg.App.HTTP.Address)
	}
	if cfg.App.HTTP.HTTPTimeout.ReadTimeout != 4*time.Second {
		t.Fatalf("read timeout want 4s, got %v", cfg.App.HTTP.HTTPTimeout.ReadTimeout)
	}
	if cfg.App.HTTP.User != "hihika" || cfg.App.HTTP.Password != "pwd" {
		t.Fatalf("auth not parsed")
	}
	if cfg.DB.MaxOpenConns != 20 || cfg.DB.MaxIdleConns != 5 {
		t.Fatalf("db pool settings mismatch")
	}
	
}

func TestMustLoad_MissingFile(t *testing.T) {
	t.Setenv("CONFIG_PATH", "./definitely_not_exists.yaml")

	if os.Getenv("BE_CHILD") == "1" {
		_ = config.MustLoad()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestMustLoad_MissingFile")
	cmd.Env = append(os.Environ(), "BE_CHILD=1")
	err := cmd.Run()
	if err == nil {
		t.Fatalf("expected MustLoad to fail on missing config file")
	}
}

func TestMustLoad_ParsesDurations(t *testing.T) {
	yaml := `
env: "local"
app:
  name: "x"
  http:
    address: "localhost:8081"
    timeouts:
      read: 7s
      write: 8s
      idle: 9s
      event: 10s
db:
  dsn: "postgres://u:p@db:5432/meetings?sslmode=disable"
  max_open_conns: 20
  max_idle_conns: 5
  conn_max_lifetime: "45m"
`
	p := writeTempConfig(t, yaml)
	t.Setenv("CONFIG_PATH", p)

	cfg := config.MustLoad()

	if cfg.App.HTTP.HTTPTimeout.ReadTimeout != 7*time.Second {
		t.Fatalf("wrong read timeout: %v", cfg.App.HTTP.HTTPTimeout.ReadTimeout)
	}
	if cfg.App.HTTP.HTTPTimeout.WriteTimeout != 8*time.Second {
		t.Fatalf("wrong write timeout: %v", cfg.App.HTTP.HTTPTimeout.WriteTimeout)
	}
	if cfg.App.HTTP.HTTPTimeout.IdleTimeout != 9*time.Second {
		t.Fatalf("wrong idle timeout: %v", cfg.App.HTTP.HTTPTimeout.IdleTimeout)
	}
	if cfg.App.HTTP.HTTPTimeout.EventTimeout != 10*time.Second {
		t.Fatalf("wrong event timeout: %v", cfg.App.HTTP.HTTPTimeout.EventTimeout)
	}
	if cfg.DB.ConnMaxLifetime != 45*time.Minute {
		t.Fatalf("wrong conn max lifetime: %v", cfg.DB.ConnMaxLifetime)
	}

}
