package domain

import "time"

type TimeStamps struct {
	InsertedAt time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
}
