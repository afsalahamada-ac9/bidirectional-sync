/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package repository

import (
	"database/sql"
	"time"

	"sudhagar/glad/entity"
)

// CoursePGSQL mysql repo
type CoursePGSQL struct {
	db *sql.DB
}

// NewCoursePGSQL create new repository
func NewCoursePGSQL(db *sql.DB) *CoursePGSQL {
	return &CoursePGSQL{
		db: db,
	}
}

// Create creates a course
func (r *CoursePGSQL) Create(e *entity.Course) (entity.ID, error) {
	stmt, err := r.db.Prepare(`
		INSERT INTO course (id, tenant_id, ext_id, center_id, name, notes, timezone, location, status, ctype, max_attendees, num_attendees, is_auto_approve, created_at) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`)
	if err != nil {
		return e.ID, err
	}
	_, err = stmt.Exec(
		e.ID,
		e.TenantID,
		e.ExtID,
		e.CenterID,
		e.Name,
		e.Notes,
		e.Timezone,
		e.Location, // TODO: to be converted into json
		int(e.Status),
		int(e.CType),
		e.MaxAttendees,
		e.NumAttendees,
		e.IsAutoApprove,
		time.Now().Format("2006-01-02"),
	)
	if err != nil {
		return e.ID, err
	}
	err = stmt.Close()
	if err != nil {
		return e.ID, err
	}
	return e.ID, nil
}

// Get retrieves a course
func (r *CoursePGSQL) Get(id entity.ID) (*entity.Course, error) {
	stmt, err := r.db.Prepare(`
		SELECT id, tenant_id, ext_id, center_id, name, notes, timezone, location,
		status, ctype, max_attendees, num_attendees, id_auto_approve, created_at
		FROM course
		WHERE id = $1;`)
	if err != nil {
		return nil, err
	}
	var c entity.Course
	var ext_id sql.NullString
	var name, notes, timezone, loc_json sql.NullString
	err = stmt.QueryRow(id).Scan(&c.ID, &c.TenantID, &ext_id, &name, &notes, &timezone, &loc_json,
		&c.Status, &c.CType, &c.MaxAttendees, &c.NumAttendees, &c.IsAutoApprove, &c.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	c.ExtID = ext_id.String
	c.Name = name.String
	c.Notes = notes.String
	c.Timezone = timezone.String
	c.Location = entity.CourseLocation{} // TODO: Convert to JSON object: loc_json.String

	return &c, nil
}

// Update updates a course
func (r *CoursePGSQL) Update(e *entity.Course) error {
	e.UpdatedAt = time.Now()
	loc_json := "" // TODO: convert location to json_str

	_, err := r.db.Exec(`
		UPDATE course SET center_id = $1, name = $2, notes = $3, timezone = $4, location = $5,
		status = $6, ctype = $7, max_attendees = $8, num_attendees = $9, is_auto_approve = $10,
		updated_at = $11
		WHERE id = $12;`,
		e.CenterID, e.Name, e.Notes, e.Timezone, loc_json, int(e.Status), int(e.CType),
		e.MaxAttendees, e.NumAttendees, e.IsAutoApprove, e.UpdatedAt.Format("2006-01-02"), e.ID)
	if err != nil {
		return err
	}
	return nil
}

// Search searches courses
func (r *CoursePGSQL) Search(tenantID entity.ID,
	query string,
) ([]*entity.Course, error) {
	stmt, err := r.db.Prepare(`
		SELECT id, tenant_id, ext_id, center_id, name, notes, timezone, location,
		status, ctype, max_attendees, num_attendees, id_auto_approve, created_at
		FROM course
		WHERE tenant_id = $1 AND name LIKE $2;`)
	if err != nil {
		return nil, err
	}
	var courses []*entity.Course
	rows, err := stmt.Query(tenantID, "%"+query+"%")
	if err != nil {
		return nil, err
	}

	var ext_id sql.NullString
	var name, notes, timezone, loc_json sql.NullString

	for rows.Next() {
		var c entity.Course
		err = rows.Scan(&c.ID, &c.TenantID, &ext_id, &name, &notes, &timezone, &loc_json,
			&c.Status, &c.CType, &c.MaxAttendees, &c.NumAttendees, &c.IsAutoApprove, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		c.ExtID = ext_id.String
		c.Name = name.String
		c.Notes = notes.String
		c.Timezone = timezone.String
		c.Location = entity.CourseLocation{} // TODO: Convert to JSON object: loc_json.String
		courses = append(courses, &c)
	}

	return courses, nil
}

// List lists courses
func (r *CoursePGSQL) List(tenantID entity.ID) ([]*entity.Course, error) {
	stmt, err := r.db.Prepare(`
		SELECT id, tenant_id, ext_id, center_id, name, notes, timezone, location,
		status, ctype, max_attendees, num_attendees, id_auto_approve, created_at
		FROM course
		WHERE tenant_id = $1;`)
	if err != nil {
		return nil, err
	}
	var courses []*entity.Course
	rows, err := stmt.Query(tenantID)
	if err != nil {
		return nil, err
	}

	var ext_id sql.NullString
	var name, notes, timezone, loc_json sql.NullString
	for rows.Next() {
		var c entity.Course
		err = rows.Scan(&c.ID, &c.TenantID, &ext_id, &name, &notes, &timezone, &loc_json,
			&c.Status, &c.CType, &c.MaxAttendees, &c.NumAttendees, &c.IsAutoApprove, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		c.ExtID = ext_id.String
		c.Name = name.String
		c.Notes = notes.String
		c.Timezone = timezone.String
		c.Location = entity.CourseLocation{} // TODO: Convert to JSON object: loc_json.String
		courses = append(courses, &c)
	}
	return courses, nil
}

// Delete deletes a course
func (r *CoursePGSQL) Delete(id entity.ID) error {
	res, err := r.db.Exec(`DELETE FROM course WHERE id = $1;`, id)
	if err != nil {
		return err
	}

	if cnt, _ := res.RowsAffected(); cnt == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Get total courses
func (r *CoursePGSQL) GetCount(tenantID entity.ID) (int, error) {
	stmt, err := r.db.Prepare(`SELECT count(*) FROM course WHERE tenant_id = $1;`)
	if err != nil {
		return 0, err
	}

	var count int
	err = stmt.QueryRow(tenantID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
