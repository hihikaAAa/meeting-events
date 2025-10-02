package main

import (
	
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	pg "github.com/hihikaAAa/meeting-events/internal/adapters/postgres"
	"github.com/hihikaAAa/meeting-events/internal/config"
	"github.com/hihikaAAa/meeting-events/internal/httpserver"
	hcreate "github.com/hihikaAAa/meeting-events/internal/httpserver/handlers/meeting/create"
	hdelete "github.com/hihikaAAa/meeting-events/internal/httpserver/handlers/meeting/delete"
	hget "github.com/hihikaAAa/meeting-events/internal/httpserver/handlers/meeting/get"
	hupdate "github.com/hihikaAAa/meeting-events/internal/httpserver/handlers/meeting/update"
	"github.com/hihikaAAa/meeting-events/internal/lib/logger/setup"
	ucancel "github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/cancel"
	ucreate "github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/create"
	uupdate "github.com/hihikaAAa/meeting-events/internal/app/usecase/meeting/update"
	mig "github.com/hihikaAAa/meeting-events/internal/adapters/postgres/migrate"
	obw "github.com/hihikaAAa/meeting-events/internal/services/outboxworker"

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
	if err := mig.Up(ctx, db.Pool); err != nil {
		log.Error("migrations failed", slog.String("err", err.Error()))
		os.Exit(1)
	}
	defer db.Pool.Close()

	uow := pg.NewUoW(db.Pool)

	if cfg.Outbox.Enabled {
		pub := obw.NewLogPublisher(log)
		worker := obw.New(log, uow, pub, cfg.Outbox.BatchSize, cfg.Outbox.PollInterval)
		worker.Start(ctx)
	}


	ucCreate := ucreate.New(uow)
	ucUpdate := uupdate.New(uow)
	ucCancel := ucancel.New(uow)

	hc := hcreate.New(log, ucCreate)
	hu := hupdate.New(log, ucUpdate)
	hd := hdelete.New(log, ucCancel)
	hg := hget.New(log, uow)
	
	handlers := httpserver.Handlers{
		Create: hc,
		Get:    hg,
		Update: hu,
		Delete: hd,
	}
	r := httpserver.NewRouter(handlers, log, cfg.App.HTTP.User, cfg.App.HTTP.Password)

	srv := &http.Server{
		Addr:         cfg.App.HTTP.Address,
		Handler:      r,
		ReadTimeout:  cfg.App.HTTP.ReadTimeout,
		WriteTimeout: cfg.App.HTTP.WriteTimeout,
		IdleTimeout:  cfg.App.HTTP.IdleTimeout,
	}

	go func() {
		log.Info("http listen", slog.String("addr", cfg.App.HTTP.Address))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("http error", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctxShut, cancelShut := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShut()
	_ = srv.Shutdown(ctxShut)
	log.Info("stopped")
}