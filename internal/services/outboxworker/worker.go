package outboxworker

import (
	"context"
	"log/slog"
	"time"

	"github.com/hihikaAAa/meeting-events/internal/app/ports"
)

type Publisher interface {
	Publish(ctx context.Context, e ports.OutboxEvent) error
}

type Worker struct {
	log        *slog.Logger
	uow        ports.UnitOfWork
	pub        Publisher
	batchSize  int
	interval   time.Duration
}

func New(log *slog.Logger, uow ports.UnitOfWork, pub Publisher, batchSize int, interval time.Duration) *Worker {
	return &Worker{
		log:       log.With(slog.String("cmp", "outbox-worker")),
		uow:       uow,
		pub:       pub,
		batchSize: batchSize,
		interval:  interval,
	}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	w.log.Info("started", slog.Duration("interval", w.interval), slog.Int("batch", w.batchSize))

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				w.log.Info("stopped")
				return
			case <-ticker.C:
				if err := w.Tick(ctx); err != nil {
					w.log.Error("tick failed", slog.String("err", err.Error()))
				}
			}
		}
	}()
}

func (w *Worker) Tick(ctx context.Context) error {
	var toMark []int64

	err := w.uow.WithinTx(ctx, func(ctx context.Context, r ports.Repos) error {
		events, err := r.Outbox().FetchPending(ctx, w.batchSize)
		if err != nil { 
			return err 
		}
		if len(events) == 0 {
			w.log.Debug("no events")
			return nil
		}

		for _, e := range events {
			if err := w.pub.Publish(ctx, e); err != nil {
				w.log.Error("publish failed",
					slog.Int64("id", e.ID),
					slog.String("type", e.EventType),
					slog.String("err", err.Error()))
				continue
			}
			toMark = append(toMark, e.ID)
		}
		if len(toMark) > 0 {
			if err := r.Outbox().MarkProcessed(ctx, toMark); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
