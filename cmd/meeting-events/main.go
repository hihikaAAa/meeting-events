package main

import (
	
	"context"
    "log/slog"
    "os"
    "time"

    "github.com/joho/godotenv"

    "github.com/hihikaAAa/meeting-events/internal/config"
    "github.com/hihikaAAa/meeting-events/internal/lib/logger/setup"
    pg "github.com/hihikaAAa/meeting-events/internal/adapters/postgres"

)

func main(){
	_ = godotenv.Load("local.env")

	cfg := config.MustLoad()

	log := setup.SetupLogger(cfg.Env)

	slog.SetDefault(log)

	ctx,cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := pg.New(
        ctx,
        cfg.DB.DSN,
        int32(cfg.DB.MaxOpenConns),       
        cfg.DB.ConnMaxLifetime,             
        5*time.Minute,                      
    )
	 if err != nil {
        log.Error("db init failed", slog.String("err", err.Error()))
        os.Exit(1)
    }
	defer db.Pool.Close()

	
	
}