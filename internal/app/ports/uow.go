package ports

import "context"


type UnitOfWork interface {
  WithinTx(ctx context.Context, fn func(ctx context.Context, r Repos) error) error
}

type Repos interface {
  Meetings() MeetingRepository
  Outbox()   OutboxRepository
}
