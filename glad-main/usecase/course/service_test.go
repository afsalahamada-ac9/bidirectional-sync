/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package course

import (
	"testing"
	"time"

	"sudhagar/glad/entity"

	"github.com/stretchr/testify/assert"
)

const (
	courseDefault entity.ID = 13790493495087071234
	tenantAlice   entity.ID = 13790492210917015554
	aliceExtID              = "000aliceExtID"
	aliceCenterID           = 13790493495087075501

	// todo: add multi-tenant support
	// tenantBob    entity.ID = 13790492210917015555
	bobExtID    = "000bobExtID"
	bobCenterID = 13790493495087075502
)

func newFixtureCourse() *entity.Course {
	return &entity.Course{
		ID:       courseDefault,
		TenantID: tenantAlice,
		ExtID:    aliceExtID,
		CenterID: aliceCenterID,
		Name:     "Course Part 1",
		Notes:    "This is a course notes. It can be multi-line text. The notes can be longer than this.",
		Timezone: "PST",
		Location: entity.CourseLocation{
			Street1: "1 Street Way",
			Street2: "",
			City:    "CityName",
			State:   "California",
			Country: "USA",
			Zip:     "12345-6789",
		},
		Status:        entity.CourseActive,
		CType:         entity.CourseInPerson,
		MaxAttendees:  50,
		NumAttendees:  12,
		IsAutoApprove: true,
		CreatedAt:     time.Now(),
	}
}

func Test_Create(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)
	tmpl := newFixtureCourse()
	_, err := m.CreateCourse(tmpl.TenantID, tmpl.ExtID, tmpl.CenterID,
		tmpl.Name, tmpl.Notes, tmpl.Timezone,
		tmpl.Location, tmpl.Status, tmpl.CType,
		tmpl.MaxAttendees, tmpl.NumAttendees, tmpl.IsAutoApprove)
	assert.Nil(t, err)
	assert.False(t, tmpl.CreatedAt.IsZero())
}

func Test_SearchAndFind(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)
	tmpl1 := newFixtureCourse()
	tmpl2 := newFixtureCourse()
	tmpl2.Name = "Course Sahaj Meditation"
	tmpl2.ExtID = bobExtID

	tID, _ := m.CreateCourse(tmpl1.TenantID, tmpl1.ExtID, tmpl1.CenterID,
		tmpl1.Name, tmpl1.Notes, tmpl1.Timezone,
		tmpl1.Location, tmpl1.Status, tmpl1.CType,
		tmpl1.MaxAttendees, tmpl1.NumAttendees, tmpl1.IsAutoApprove)
	_, _ = m.CreateCourse(tmpl2.TenantID, tmpl2.ExtID, tmpl2.CenterID,
		tmpl2.Name, tmpl1.Notes, tmpl1.Timezone,
		tmpl1.Location, tmpl1.Status, tmpl1.CType,
		tmpl1.MaxAttendees, tmpl1.NumAttendees, tmpl1.IsAutoApprove)

	t.Run("search", func(t *testing.T) {
		res, err := m.SearchCourses(tmpl1.TenantID, "Part")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(res))
		assert.Equal(t, tmpl1.ExtID, res[0].ExtID)
		assert.Equal(t, tmpl1.CenterID, res[0].CenterID)
		assert.Equal(t, tmpl1.Status, res[0].Status)
		// TODO: checks for other fields to be added

		// 'default' query value matches both the course names
		res, err = m.SearchCourses(tmpl1.TenantID, "Sahaj")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(res))

		res, err = m.SearchCourses(tmpl1.TenantID, "non-existent")
		assert.Equal(t, entity.ErrNotFound, err)
		assert.Nil(t, res)
	})
	t.Run("list all", func(t *testing.T) {
		all, err := m.ListCourses(tmpl1.TenantID)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(all))
	})

	t.Run("get", func(t *testing.T) {
		saved, err := m.GetCourse(tID)
		assert.Nil(t, err)
		assert.Equal(t, tmpl1.TenantID, saved.TenantID)
		assert.Equal(t, tmpl1.ExtID, saved.ExtID)
		assert.Equal(t, tmpl1.CenterID, saved.CenterID)
		assert.Equal(t, tmpl1.Status, saved.Status)
		assert.Equal(t, tmpl1.Name, saved.Name)
	})
}

func Test_Update(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)
	tmpl := newFixtureCourse()
	id, err := m.CreateCourse(tmpl.TenantID, tmpl.ExtID, tmpl.CenterID,
		tmpl.Name, tmpl.Notes, tmpl.Timezone,
		tmpl.Location, tmpl.Status, tmpl.CType,
		tmpl.MaxAttendees, tmpl.NumAttendees, tmpl.IsAutoApprove)

	assert.Nil(t, err)

	saved, _ := m.GetCourse(id)
	saved.CType = entity.CourseOnline
	assert.Nil(t, m.UpdateCourse(saved))

	updated, err := m.GetCourse(id)
	assert.Nil(t, err)
	assert.Equal(t, entity.CourseOnline, updated.CType)
}

func TestDelete(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)

	tmpl1 := newFixtureCourse()
	tmpl2 := newFixtureCourse()
	tmpl2.ExtID = bobExtID
	t2ID, _ := m.CreateCourse(tmpl2.TenantID, tmpl2.ExtID, tmpl2.CenterID,
		tmpl2.Name, tmpl1.Notes, tmpl1.Timezone,
		tmpl1.Location, tmpl1.Status, tmpl1.CType,
		tmpl1.MaxAttendees, tmpl1.NumAttendees, tmpl1.IsAutoApprove)

	err := m.DeleteCourse(tmpl1.ID)
	assert.Equal(t, entity.ErrNotFound, err)

	err = m.DeleteCourse(t2ID)
	assert.Nil(t, err)
	_, err = m.GetCourse(t2ID)
	assert.Equal(t, entity.ErrNotFound, err)
}
