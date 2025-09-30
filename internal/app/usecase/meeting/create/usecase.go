package create

import (
	"context"


	"github.com/google/uuid"
	app "github.com/hihikaAAa/meeting-events/internal/app/app_errors"
	"github.com/hihikaAAa/meeting-events/internal/app/ports"
	"github.com/hihikaAAa/meeting-events/internal/domain/domErrors"
	"github.com/hihikaAAa/meeting-events/internal/domain/meeting"
)

type UseCase struct{
	uow ports.UnitOfWork
}

func New(uow ports.UnitOfWork) *UseCase{
	return &UseCase{uow:uow}
}

func (uc *UseCase) Handle(ctx context.Context, in Input) (Output, error){
	m, err := meeting.NewMeeting(in.Title,in.StartsAt, in.Duration)
	if err != nil{
		switch err{
		case domErrors.ErrInvalidTitle, domErrors.ErrInvalidDuration,domErrors.ErrInvalidDuration:
				return Output{}, app.ErrValidation
		default:
			return Output{},err	
		}
	}

	var id uuid.UUID
	err = uc.uow.WithinTx(ctx, func(ctx context.Context, r ports.Repos) error{
		if err := r.Meetings().Create(ctx,m); err != nil{
			return err
		}
		id = m.ID
		if err := r.Outbox().Add(ctx, "meeting",m.ID,"MeetingCreated",map[string]any{
			"id":       m.ID.String(),
			"title":    m.Title,
			"startsAt": m.StartsAt,
		}); err!=nil{
			return err
		}
		return nil
	})
	if err != nil{
		return Output{},err
	}
	return Output{ID: id.String()}, nil
}