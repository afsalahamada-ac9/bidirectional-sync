package main

import (
	entity "csv-detector/entity"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"encoding/csv"

	"github.com/fsnotify/fsnotify"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// prepare the watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("error creating the watcher", err)
	}
	watch_directory := "./test"
	err = watcher.Add(watch_directory)
	if err != nil {
		log.Println("error adding the directory to watcher", err)
	}

	// connect the database and migrate the db
	connection_string := "host=localhost user=postgres password=1234 port=5432 dbname=democsv sslmode=disable"
	db, err := gorm.Open(postgres.Open(connection_string), &gorm.Config{})
	if err != nil {
		log.Println("there was an error connecting to the database", err)
	}
	log.Println("db connected")
	err = db.AutoMigrate(&entity.Account{}, &entity.Center{}, &entity.Course{}, &entity.CourseContact{}, &entity.CourseNotify{}, &entity.CourseOrganizer{}, &entity.CourseTeacher{}, &entity.CourseTiming{}, &entity.Tenant{})
	if err != nil {
		log.Println("there was an error migrating the database", err)
	}
	log.Println("migration successful")
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Create == fsnotify.Create && strings.HasSuffix(event.Name, ".csv") {
				log.Println("detected new csv file", event.Name)
				processCsv(event.Name, db)
			}
		case err := <-watcher.Errors:
			log.Println("error occurred", err)
		}
	}
}
func processCsv(filePath string, db *gorm.DB) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("error navigating to the filepath", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Println("reader error:", err)
		return
	}
	if len(records) < 2 {
		log.Println("empty csv file detected")
		return
	}

	header := records[0]
	var entityType string

	if contains(header, "name") && contains(header, "country") {
		entityType = "tenant"
	} else if contains(header, "center_name") {
		entityType = "center"
	} else if contains(header, "course_name") {
		entityType = "course"
	} else if contains(header, "username") {
		entityType = "account"
	} else if contains(header, "course_id") && contains(header, "organizer_id") {
		entityType = "course_organizer"
	} // Add more conditions for other entities
	for _, row := range records[1:] { // Skip header row
		switch entityType {
		case "tenant":
			// Assuming tenant.csv has the correct columns
			tenant := entity.Tenant{
				Name:      row[1],
				Country:   row[2],
				CreatedAt: parseTime(row[3]),
				UpdatedAt: parseTime(row[4]),
			}
			db.Create(&tenant)

		case "center":
			center := entity.Center{
				ExtId:            row[1],
				CenterName:       row[2],
				Location:         entity.JSONB(row[3]), // Assuming JSONB is a string
				GeoLocation:      entity.JSONB(row[4]),
				Capacity:         parseInt(row[5]),
				Mode:             row[6],
				IsNationalCenter: parseBool(row[7]),
				CreatedAt:        parseTime(row[8]),
				UpdatedAt:        parseTime(row[9]),
			}
			db.Create(&center)

		case "course":
			course := entity.Course{
				ExtId:         row[1],
				Name:          row[2],
				Notes:         row[3],
				Status:        row[4],
				MaxAttendees:  parseInt(row[5]),
				TimeZone:      row[6],
				Location:      entity.JSONB(row[7]),
				CenterId:      parseInt(row[8]),
				CType:         row[9],
				NumAttendees:  parseInt(row[10]),
				IsAutoApprove: parseBool(row[11]),
				CreatedAt:     parseTime(row[12]),
				UpdatedAt:     parseTime(row[13]),
			}
			db.Create(&course)

		case "account":
			account := entity.Account{
				ExtId:     row[1],
				Username:  row[2],
				FirstName: row[3],
				LastName:  row[4],
				Phone:     row[5],
				Email:     row[6],
				Type:      row[7],
				CreatedAt: parseTime(row[8]),
				UpdatedAt: parseTime(row[9]),
			}
			db.Create(&account)

		case "course_organizer":
			courseOrganizer := entity.CourseOrganizer{
				CourseId:    parseInt(row[1]),
				OrganizerId: parseInt(row[2]),
				UpdatedAt:   parseTime(row[3]),
			}
			db.Create(&courseOrganizer)

		case "course_contact":
			courseContact := entity.CourseContact{
				CourseId:  parseInt(row[1]),
				ContactId: parseInt(row[2]),
				UpdatedAt: parseTime(row[3]),
			}
			db.Create(&courseContact)

		case "course_teacher":
			courseTeacher := entity.CourseTeacher{
				CourseId:  parseInt(row[1]),
				TeacherId: parseInt(row[2]),
				UpdatedAt: parseTime(row[3]),
			}
			db.Create(&courseTeacher)

		case "course_timing":
			courseTiming := entity.CourseTiming{
				CourseId:   parseInt(row[1]),
				ExtId:      row[2],
				CourseDate: parseTime(row[3]),
				StartTime:  parseTime(row[4]),
				EndTime:    parseTime(row[5]),
				UpdatedAt:  parseTime(row[6]),
			}
			db.Create(&courseTiming)

		case "course_notify":
			courseNotify := entity.CourseNotify{
				CourseId:  parseInt(row[1]),
				NotifyId:  parseInt(row[2]),
				UpdatedAt: parseTime(row[3]),
			}
			db.Create(&courseNotify)

		default:
			log.Println("unknown entity type:", entityType)
		}
	}
}

// Helper functions to parse data types
func parseTime(value string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", value)
	return t
}

func parseInt(value string) int {
	i, _ := strconv.Atoi(value)
	return i
}

func parseBool(value string) bool {
	b, _ := strconv.ParseBool(value)
	return b
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
