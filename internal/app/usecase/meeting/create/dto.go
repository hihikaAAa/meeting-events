package create 

import "time"

type Input struct{
	Title string
	StartsAt time.Time
	Duration time.Duration
	OwnerID string
}

type Output struct{
	ID string
}