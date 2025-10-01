package meeting

import (
	"time"

	"github.com/google/uuid"
	"github.com/hihikaAAa/meeting-events/internal/domain/domErrors"
)

const (
  minDuration = time.Minute
  maxDuration = 480 * time.Minute
  maxTitleLen = 200
)

type Meeting struct{
	ID uuid.UUID
	Title string
	StartsAt time.Time
	Duration  time.Duration
	Status Status
	createdAt time.Time
	updatedAt time.Time
	events []any
}

type Status string

const (
    StatusScheduled Status = "scheduled"
    StatusCanceled  Status = "canceled"
	StatusOngoing Status = "ongoing"
)

func NewMeeting (title string, startsAt time.Time, duration time.Duration)(*Meeting, error){
	if len(title)==0 || len(title) > maxTitleLen{
		return nil, domErrors.ErrInvalidTitle
	}
	if startsAt.Before(time.Now()){
		return nil, domErrors.ErrInvalidTime
	}
	if duration < minDuration|| duration > maxDuration{
		return nil, domErrors.ErrInvalidDuration
	}

	m := &Meeting{
		ID : uuid.New(),
		Title: title,
		StartsAt: startsAt,
		Duration: duration,
		Status: StatusScheduled,
		createdAt: time.Now(),
		updatedAt: time.Now(),
		
	}
	m.addEvent(MeetingCreated{ID: m.ID})
	return m, nil
}

func (m *Meeting) Cancel() error{
	if m.Status == StatusCanceled{
		return domErrors.ErrAlreadyCanceled
	}
	if m.Status ==  StatusOngoing{
		return domErrors.ErrOngoing
	}
	m.Status = StatusCanceled
	m.updatedAt = time.Now()
	m.addEvent(MeetingCanceled{ID: m.ID})
	return nil
}

func (m *Meeting) Update(title string, startsAt time.Time, duration time.Duration) error{
	if m.Status == StatusCanceled{
		return domErrors.ErrAlreadyCanceled
	}
	if m.Status == StatusOngoing{
		return domErrors.ErrOngoing
	}

	now := time.Now()
    changed := false

	if title != "" {
        if l := len(title); l < 1 || l > maxTitleLen {
            return domErrors.ErrInvalidTitle
        }
    }
	if !startsAt.IsZero() {
        if !startsAt.After(now) {
            return domErrors.ErrInvalidTime
        }
    }
	if duration != 0 {
        if duration < minDuration || duration > maxDuration {
            return domErrors.ErrInvalidDuration
        }
    }

	if title != "" && title != m.Title {
        m.Title = title
        changed = true
    }
    if !startsAt.IsZero() && !startsAt.Equal(m.StartsAt) {
        m.StartsAt = startsAt
        changed = true
    }
    if duration != 0 && duration != m.Duration {
        m.Duration = duration
        changed = true
    }

	if !changed{
		return nil
	}
	
	m.updatedAt = time.Now()
	m.addEvent(MeetingUpdated{ID: m.ID})
	return nil
}

func (m *Meeting) Events() []any {
    return m.events
}
func (m *Meeting) addEvent(e any) {
    m.events = append(m.events, e)
}
