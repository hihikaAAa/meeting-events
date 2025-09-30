package tools

import (
	"context"

	"github.com/google/uuid"

	"github.com/hihikaAAa/meeting-events/internal/app/ports"
	"github.com/hihikaAAa/meeting-events/internal/domain/meeting"
)

type FakeRepos struct {
	Mr ports.MeetingRepository
	Or ports.OutboxRepository
}

func (f FakeRepos) Meetings() ports.MeetingRepository { 
	return f.Mr
}
func (f FakeRepos) Outbox() ports.OutboxRepository    { 
	return f.Or 
}

type FakeUoW struct {
	Repos FakeRepos
	Err   error
}

func (u FakeUoW) WithinTx(ctx context.Context, fn func(ctx context.Context, r ports.Repos) error) error {
	if u.Err != nil {
		return u.Err
	}
	return fn(ctx, u.Repos)
}

type MockMeetRepo struct {
	Created *meeting.Meeting
}

func (m *MockMeetRepo) Create(ctx context.Context, mm *meeting.Meeting) error {
	m.Created = mm
	return nil
}
func (m *MockMeetRepo) GetByID(context.Context, uuid.UUID) (*meeting.Meeting, error) { 
	return nil, nil 
}
func (m *MockMeetRepo) Update(context.Context, *meeting.Meeting) error{
	 return nil 
}
func (m *MockMeetRepo) Cancel(context.Context, uuid.UUID) error{ 
	return nil
}

type MockOutbox struct {
	Added int
}

func (m *MockOutbox) Add(context.Context, string, uuid.UUID, string, any) error {
	m.Added++
	return nil
}
