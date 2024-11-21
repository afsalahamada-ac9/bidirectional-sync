package entity

import "time"

type CourseTeacher struct {
	CourseId  int
	TeacherId int
	UpdatedAt time.Time
}
