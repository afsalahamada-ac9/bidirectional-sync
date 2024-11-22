package entity

import "time"

type CourseNotify struct {
	CourseId  int
	NotifyId  int
	UpdatedAt time.Time
}
