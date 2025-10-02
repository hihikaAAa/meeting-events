package migrate

import (
    "context"
    "embed"
    "fmt"
    "sort"
    "strings"

    "github.com/jackc/pgx/v5/pgxpool"
)

var files embed.FS

func Up(ctx context.Context, pool *pgxpool.Pool) error {
    entries, err := files.ReadDir("../init")
    if err != nil { 
		return err
	}

    names := make([]string, 0, len(entries))
    for _, e := range entries {
        if e.IsDir() { 
			continue 
		}
        if strings.HasSuffix(e.Name(), ".sql") { 
			names = append(names, e.Name()) 
		}
    }
    sort.Strings(names)

    for _, name := range names {
        b, err := files.ReadFile("../init/" + name)
        if err != nil { return err }
        if _, err := pool.Exec(ctx, string(b)); err != nil {
            return fmt.Errorf("apply %s: %w", name, err)
        }
    }
    return nil
}
