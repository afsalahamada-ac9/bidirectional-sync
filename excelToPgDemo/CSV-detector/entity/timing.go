package entity

import "time"

type CourseTiming struct {
	CourseId   int
	ExtId      string
	CourseDate time.Time
	StartTime  time.Time
	EndTime    time.Time
	UpdatedAt  time.Time
}
