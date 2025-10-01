package cancel

import (
	"context"


	"github.com/google/uuid"
	app "github.com/hihikaAAa/meeting-events/internal/app/app_errors"
	"github.com/hihikaAAa/meeting-events/internal/app/ports"
)

type UseCase struct{
	uow ports.UnitOfWork
}

func New(uow ports.UnitOfWork) *UseCase{
	return &UseCase{uow:uow}
}

func (uc *UseCase) Handle(ctx context.Context, in Input)(Output, error){
	id, err := uuid.Parse(in.ID)
	if err != nil{
		return Output{}, app.ErrValidation
	}

	var out Output
	err = uc.uow.WithinTx(ctx, func(ctx context.Context, r ports.Repos) error{
		m,err := r.Meetings().GetByID(ctx,id)
		if err !=nil{
			return app.ErrNotFound
		}
		if err = m.Cancel(); err != nil{
			return app.ErrConflict
		}
		if err = r.Meetings().Update(ctx,m); err != nil{
			return err
		}
		if err:= r.Outbox().Add(ctx, "meeting", m.ID, "MeetingCanceled", map[string]any{"id": m.ID.String()}); err !=nil{
			return err
		}
		out = Output{ID: m.ID.String()}
		return nil
	})
	if err != nil{
		return Output{}, err
	}
	return out,nil

}