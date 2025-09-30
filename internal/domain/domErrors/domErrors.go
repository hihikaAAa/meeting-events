package domErrors 

import(
	"errors"
)

var(
	ErrInvalidTitle = errors.New("invalid title")
	ErrInvalidTime = errors.New("startsAt must be in the future")
	ErrInvalidDuration = errors.New("invalid meeting duration")
	ErrAlreadyCanceled = errors.New("meeting has been already canceled")
	ErrOngoing = errors.New("Ongoing meeting cannot be canceled")
)