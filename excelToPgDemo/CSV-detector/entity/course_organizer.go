package entity

import "time"

type CourseOrganizer struct {
	CourseId    int
	OrganizerId int
	UpdatedAt   time.Time
}
