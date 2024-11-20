package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Models
type JSONB string

type Tenant struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique;not null"`
	Country   string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Course struct {
	ID           uint   `gorm:"primaryKey"`
	ExtID        string `gorm:"unique;not null"`
	Name         string `gorm:"not null"`
	Notes        string
	Status       string `gorm:"not null;default:'draft'"`
	MaxAttendees int
	Timezone     string
	Location     JSONB
	CenterID     uint   `gorm:"not null"`
	Type         string `gorm:"not null;default:'in-person'"`
	NumAttendees int    `gorm:"default:0"`
	AutoApprove  bool   `gorm:"default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Account struct {
	ID        uint   `gorm:"primaryKey"`
	ExtID     string `gorm:"unique;not null"`
	Username  string `gorm:"not null"`
	FirstName string
	LastName  string
	Phone     string
	Email     string
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Center struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	ExtID       string `gorm:"not null"`
	CenterName  string `gorm:"not null"`
	Location    JSONB
	GeoLocation JSONB
	Capacity    int
	Mode        string `gorm:"default:'in-person'"`
	Webpage     string
	IsNational  bool `gorm:"default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Note: This is the viper function
func getValue(key string) string {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("error fetching the .env file")
	}
	value, ok := viper.Get(key).(string)
	if !ok {
		log.Println("invalid key was provided")
	}
	return value
}

// Database connection
func ConnectDB() *gorm.DB {
	dsn := getValue("CONNECTION_STRING")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Set logger to log SQL queries
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Configure connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	sqlDB.SetMaxIdleConns(10) //Note: Adjust based on needs
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Migrate the schema
	err = db.AutoMigrate(&Tenant{}, &Course{}, &Account{}, &Center{}) // Add other models
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	return db
}

// Salesforce API Client
type SalesforceClient struct {
	BaseURL   string
	AuthToken string
}

func (client *SalesforceClient) CreateCourse(course Course) error {
	url := getValue("AOL_EVENT")
	body, _ := json.Marshal(course)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+client.AuthToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create course: %s", resp.Status)
	}
	return nil
}

// Synchronization Logic
var isSyncing bool

func SyncCourses() {
	if isSyncing {
		return // Prevent infinite callback
	}
	isSyncing = true
	defer func() { isSyncing = false }()

	database := ConnectDB()
	var courses []Course
	database.Find(&courses)

	client := SalesforceClient{BaseURL: getValue("AOL_EVENT"), AuthToken: getValue("AUTH_TOKEN")}

	for _, course := range courses {
		err := client.CreateCourse(course)
		if err != nil {
			log.Println("Error syncing course to Salesforce:", err)
		}
	}
}
func InsertSampleData(db *gorm.DB) {
	tenants := []Tenant{
		{Name: "Tenant A", Country: "USA"},
		{Name: "Tenant B", Country: "Canada"},
	}

	courses := []Course{
		{ExtID: "C001", Name: "Course 1", Status: "draft", CenterID: 1},
		{ExtID: "C002", Name: "Course 2", Status: "draft", CenterID: 1},
	}

	accounts := []Account{
		{ExtID: "A001", Username: "user1", FirstName: "John", LastName: "Doe", Email: "john@example.com"},
		{ExtID: "A002", Username: "user2", FirstName: "Jane", LastName: "Doe", Email: "jane@example.com"},
	}

	centers := []Center{
		{ExtID: "C001", CenterName: "Center 1", Capacity: 100},
		{ExtID: "C002", CenterName: "Center 2", Capacity: 200},
	}
	// Note: Here we're doing upsert
	db.Save(&tenants)
	db.Save(&courses)
	db.Save(&accounts)
	db.Save(&centers)
}

func main() {
	database := ConnectDB()
	conn, err := pgx.Connect(context.Background(), getValue("CONNECTION_STRING"))
	if err != nil {
		log.Println("there was an error connecting to the database")
	}
	_, err = conn.Exec(context.Background(), "LISTEN events")
	if err != nil {
		log.Println("error listening to the channel")
	}
	log.Println("listening successful")
	channel := make(chan *pgconn.Notification)
	go func() {
		log.Println("listener active")
		for {
			notification, err := conn.WaitForNotification(context.Background())
			if err != nil {
				log.Println("there was an error listening to the notification")
				continue
			}
			channel <- notification
		}
	}()
	InsertSampleData(database)
	// Start the synchronization process
	SyncCourses()
	for notification := range channel {
		log.Println("received notification:", notification.Payload)
	}
}
