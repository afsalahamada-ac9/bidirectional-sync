package entity

import "time"

type CourseContact struct {
	CourseId  int
	ContactId int
	UpdatedAt time.Time
}
