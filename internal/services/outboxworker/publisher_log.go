package outboxworker

import (
	"context"
	"log/slog"

	"github.com/hihikaAAa/meeting-events/internal/app/ports"
)

type LogPublisher struct { log *slog.Logger }

func NewLogPublisher(log *slog.Logger) *LogPublisher { return &LogPublisher{log: log} }

func (p *LogPublisher) Publish(ctx context.Context, e ports.OutboxEvent) error {
	p.log.Info("publish",
		slog.Int64("id", e.ID),
		slog.String("type", e.EventType),
		slog.String("aggregate", e.Aggregate),
		slog.String("aggregate_id", e.AggregateID.String()),
	)
	return nil
}
