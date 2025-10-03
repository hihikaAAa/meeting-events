package migrate

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
    

	"github.com/jackc/pgx/v5/pgxpool"

)
func Up(ctx context.Context, pool *pgxpool.Pool, dir string) error {
	dir = strings.TrimPrefix(dir, "file://")

	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	info, err := os.Stat(abs)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return &fs.PathError{Op: "open", Path: abs, Err: fs.ErrNotExist}
	}

	entries, err := os.ReadDir(abs)
	if err != nil {
		return err
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToLower(e.Name()), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, name := range files {
		full := filepath.Join(abs, name)
		b, err := os.ReadFile(full)
		if err != nil {
			return err
		}
		if _, err := pool.Exec(ctx, string(b)); err != nil {
			return &fs.PathError{Op: "exec", Path: full, Err: err}
		}
	}
	return nil
}
