package postgres

import(
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct{
	Pool *pgxpool.Pool
}

func New(ctx context.Context, dsn string, maxConns int32, maxConnLifetime, maxConnIdle time.Duration) (*DB, error){
	cfg, err:= pgxpool.ParseConfig(dsn)
	if err != nil{
		return nil,err
	}
	cfg.MaxConns = maxConns
    cfg.MaxConnLifetime = maxConnLifetime
    cfg.MaxConnIdleTime = maxConnIdle

	pool, err := pgxpool.NewWithConfig(ctx,cfg)
	if err != nil{
		return nil, err
	}
	ctxPing, cancel := context.WithTimeout(ctx,5*time.Second)
	defer cancel()
	if err := pool.Ping(ctxPing); err != nil { 
		return nil, err 
	}
  	return &DB{Pool: pool}, nil
}