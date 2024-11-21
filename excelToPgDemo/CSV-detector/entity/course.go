package entity

import "time"

type Course struct {
	Id            int
	ExtId         string
	Name          string
	Notes         string // VARCHAR(1024)
	Status        string
	MaxAttendees  int // INTEGER
	TimeZone      string
	Location      JSONB // JSONB
	CenterId      int
	CType         string
	NumAttendees  int
	IsAutoApprove bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
