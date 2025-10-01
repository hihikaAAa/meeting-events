package update 

import(
	"time"
)
type Input struct{
	ID string
	Title string
	StartsAt time.Time
	Duration time.Duration
}

type Output struct{
	ID string
}