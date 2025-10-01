package meeting

import(
	"github.com/google/uuid"
)

type MeetingCreated struct {
    ID uuid.UUID
}

type MeetingCanceled struct {
    ID uuid.UUID
}

type MeetingUpdated struct{
    ID uuid.UUID
}
