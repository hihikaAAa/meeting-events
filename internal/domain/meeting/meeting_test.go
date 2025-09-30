package meeting_test

import (
	"testing"
	"time"

	"github.com/hihikaAAa/meeting-events/internal/domain/domErrors"
	"github.com/hihikaAAa/meeting-events/internal/domain/meeting"
	"github.com/stretchr/testify/require"
)

func TestMeeting_Create_Success(t *testing.T) {
	m, err := meeting.NewMeeting("Demo", time.Now().Add(time.Hour), 90*time.Minute)
	require.NoError(t, err)
	require.Equal(t, "Demo", m.Title)
	require.Equal(t, meeting.StatusScheduled, m.Status)
	require.Len(t, m.Events(), 1)
}

func TestMeeting_Create_InvalidTitle(t *testing.T) {
	_, err := meeting.NewMeeting("", time.Now().Add(time.Hour), 90*time.Minute)
	require.Equal(t, err, domErrors.ErrInvalidTitle)
}

func TestMeeting_Create_PastTime(t *testing.T) {
	_, err := meeting.NewMeeting("Demo", time.Now().Add(-time.Hour), 90*time.Minute)
	require.Equal(t, err, domErrors.ErrInvalidTime)
}

func TestMeeting_Create_InvalidDuration1(t *testing.T){
	_, err := meeting.NewMeeting("Demo", time.Now().Add(time.Hour), 481*time.Minute)
	require.Equal(t, err, domErrors.ErrInvalidDuration)
}

func TestMeeting_Create_InvalidDuration2(t *testing.T){
	_, err := meeting.NewMeeting("Demo", time.Now().Add(time.Hour), 30*time.Second)
	require.Equal(t, err, domErrors.ErrInvalidDuration)
}