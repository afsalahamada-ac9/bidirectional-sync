/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package entity

import (
	"time"
)

// Course type
type CourseType int

const (
	CourseInPerson CourseType = iota
	CourseOnline
	// Add new types here
)

// Course status
type CourseStatus int

const (
	CourseDraft CourseStatus = iota
	CourseArchived
	CourseOpen
	CourseExpenseSubmitted
	CourseExpenseDeclined
	CourseClosed
	CourseActive
	CourseDeclined
	CourseSubmitted
	CourseCanceled
	CoursedInactive
	// Add new types here
)

// Course Location
type CourseLocation struct {
	Street1 string
	Street2 string
	City    string
	State   string
	Zip     string
	Country string
}

// Course data
type Course struct {
	ID       ID
	TenantID ID
	CenterID ID

	ExtID string

	Name     string
	Notes    string
	Timezone string

	Location CourseLocation

	Status CourseStatus

	// TODO: CType to be renamed to delivery mode or method
	CType CourseType

	MaxAttendees int32
	NumAttendees int32

	// TODO: AutoApprove may not be required here. It's likely in Course Master or Catalog.
	IsAutoApprove bool

	// meta data
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewCourseLocation creates a new course location
func NewCourseLocation(street1 string,
	street2 string,
	city string,
	state string,
	zip string,
	country string) (*CourseLocation, error) {

	l := &CourseLocation{
		Street1: street1,
		Street2: street2,
		City:    city,
		State:   state,
		Zip:     zip,
		Country: country,
	}
	err := l.Validate()
	if err != nil {
		return nil, ErrInvalidEntity
	}
	return l, nil
}

// Validate validates course location
func (l *CourseLocation) Validate() error {
	if l.Street1 == "" || l.City == "" || l.State == "" || l.Zip == "" || l.Country == "" {
		return ErrInvalidEntity
	}
	return nil
}

// NewCourse create a new course
func NewCourse(tenantID ID,
	extID string,
	centerID ID,
	name string,
	notes string,
	timezone string,
	location CourseLocation,
	status CourseStatus,
	ctype CourseType,
	maxAttendees int32,
	numAttendees int32,
	isAutoApprove bool) (*Course, error) {
	c := &Course{
		ID:            NewID(),
		TenantID:      tenantID,
		ExtID:         extID,
		CenterID:      centerID,
		Name:          name,
		Notes:         notes,
		Timezone:      timezone,
		Location:      location,
		Status:        status,
		CType:         ctype,
		MaxAttendees:  maxAttendees,
		NumAttendees:  numAttendees,
		IsAutoApprove: isAutoApprove,
		CreatedAt:     time.Now(),
	}
	err := c.Validate()
	if err != nil {
		return nil, ErrInvalidEntity
	}
	return c, nil
}

// Validate validate course
func (c *Course) Validate() error {
	if c.ExtID == "" || c.Name == "" {
		return ErrInvalidEntity
	}
	return nil
}
