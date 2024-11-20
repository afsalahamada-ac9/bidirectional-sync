package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Tenant struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Course struct {
	ID           int64     `json:"id"`
	ExtID        string    `json:"ext_id"`
	Name         string    `json:"name"`
	Notes        string    `json:"notes"`
	Status       string    `json:"status"` // course_status
	MaxAttendees int       `json:"max_attendees"`
	Timezone     string    `json:"timezone"` // timezone_type
	Location     string    `json:"location"` // JSONB
	CenterID     int64     `json:"center_id"`
	Type         string    `json:"type"` // course_type
	NumAttendees int       `json:"num_attendees"`
	AutoApprove  bool      `json:"auto_approve"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CourseOrganizer struct {
	CourseID    int64     `json:"course_id"`
	OrganizerID int64     `json:"organizer_id"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CourseContact struct {
	CourseID  int64     `json:"course_id"`
	ContactID int64     `json:"contact_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CourseTeacher struct {
	CourseID  int64     `json:"course_id"`
	TeacherID int64     `json:"teacher_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CourseTiming struct {
	CourseID   int64     `json:"course_id"`
	ExtID      string    `json:"ext_id"`
	CourseDate time.Time `json:"course_date"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CourseNotify struct {
	CourseID  int64     `json:"course_id"`
	NotifyID  int64     `json:"notify_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Account struct {
	ID        int64     `json:"id"`
	ExtID     string    `json:"ext_id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Type      string    `json:"type"` // account_type
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FIXME: jsonb should be kept as an object
type Center struct {
	ID               int64     `json:"id"`
	ExtID            string    `json:"ext_id"`
	CenterName       string    `json:"center_name"`
	Location         string    `json:"location"`     // JSONB
	GeoLocation      string    `json:"geo_location"` // JSONB
	Capacity         int       `json:"capacity"`
	Mode             string    `json:"mode"` // center_mode
	Webpage          string    `json:"webpage"`
	IsNationalCenter bool      `json:"is_national_center"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CenterContact struct {
	CenterID  int64     `json:"center_id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updated_at"`
}

func getValue(key string) string {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("invalid configuration file error")
	}
	value, ok := viper.Get(key).(string)
	if !ok {
		log.Println("error fetching the specified key")
	}
	return value
}

func main() {
	connectionString := getValue("CONNECTION_STRING")
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Println("an error occurred when connecting to the database")
	}
	// db.AutoMigrate(&Course{}, &CourseTiming{}, &CourseTeacher{}, &CourseOrganizer{}, &CourseContact{}, &CourseNotify{}, &Tenant{}, &Center{}, &CenterContact{}, &Account{})
	err = db.AutoMigrate(&Course{})
	if err != nil {
		log.Println("there was an error migrating to the schema")
	}
	conn, err := pgx.Connect(context.Background(), getValue("CONNECTION_STRING"))
	if err != nil {
		log.Println("there was an error connecting to the daatbase")
	}
	_, err = conn.Exec(context.Background(), "LISTEN events")
	if err != nil {
		log.Println("listener service failed")
	}
	channel := make(chan *pgconn.Notification)
	go func() {
		for {
			notification, err := conn.WaitForNotification(context.Background())
			if err != nil {
				log.Println("error listening to the notification")
				continue
			}
			channel <- notification
		}
	}()
	mockData := Course{
		ID:           1,
		ExtID:        "ext-001",
		Name:         "Introduction to Go",
		Notes:        "This course covers the basics of Go programming.",
		Status:       "active",
		MaxAttendees: 30,
		Timezone:     "UTC",
		Location:     "{\"address\": \"123 Main St\", \"city\": \"Anytown\", \"state\": \"CA\"}",
		CenterID:     1,
		Type:         "online",
		NumAttendees: 10,
		AutoApprove:  true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	db.Save(&mockData)
	for msg := range channel {
		log.Println("receive notification:", msg)
	}
}

// Copyright 2024 AboveCloud9.AI Products and Services Private Limited
// All rights reserved.
// This code may not be used, copied, modified, or distributed without explicit permission.
