package main

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fsnotify/fsnotify"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Record struct {
	Id                   string   `csv:"Id"`
	Name                 string   `csv:"Name"`
	StreetAddress1       string   `csv:"Street_Address_1__c"`
	StreetAddress2       *string  `csv:"Street_Address_2__c"`
	City                 *string  `csv:"City__c"`
	State                *string  `csv:"State__c"`
	PostalOrZipCode      *string  `csv:"Postal_Or_Zip_Code__c"`
	GeolocationLatitude  *float64 `csv:"Geolocation__c.latitude"`
	GeolocationLongitude *float64 `csv:"Geolocation__c.longitude"`
	MaxCapacity          *int     `csv:"Max_Capacity__c"`
	CenterMode           *string  `csv:"Center_Mode__c"`
	IsNationalCenter     bool     `csv:"Is_National_Center__c"`
	CenterURL            string   `csv:"Center_URL__c"`
}

func main() {
	url := "host=localhost user=postgres password=1234 dbname=democsv port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Println("error connecting with the database", err)
		return
	}
	err = db.AutoMigrate(&Record{})
	if err != nil {
		log.Println("migration failed", err)
	}
	log.Println("migration success")
	watchDir := "./csv_files" // directory for new csv files
	watchForNewFiles(watchDir, db)
}

func watchForNewFiles(watchDir string, db *gorm.DB) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("watcher service failed", err)
		return
	}
	defer watcher.Close()
	err = watcher.Add(watchDir)
	if err != nil {
		log.Println("error with watcher service", err)
		return
	}
	log.Println("watching directory:", watchDir)
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				if filepath.Ext(event.Name) == ".csv" {
					log.Println("New csv file has been detected:", event.Name)
					processCSVFile(event.Name, db)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error occurred:", err)
		}
	}
}

func processCSVFile(filePath string, db *gorm.DB) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("an error occurred when opening the file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Println("failed to read the csv file:", err)
		return
	}

	// Skip the header row
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}

		var streetAddress2 *string
		if record[3] != "" {
			streetAddress2 = &record[3]
		}

		var city *string
		if record[4] != "" {
			city = &record[4]
		}

		var state *string
		if record[5] != "" {
			state = &record[5]
		}

		var postalOrZipCode *string
		if record[6] != "" {
			postalOrZipCode = &record[6]
		}

		var maxCapacity *int
		if record[9] != "" {
			capacity, _ := strconv.Atoi(record[9])
			maxCapacity = &capacity
		}

		var centerMode *string
		if record[10] != "" {
			centerMode = &record[10]
		}

		latitude, _ := strconv.ParseFloat(record[7], 64)
		longitude, _ := strconv.ParseFloat(record[8], 64)

		newRecord := Record{
			Id:                   record[0],
			Name:                 record[1],
			StreetAddress1:       record[2],
			StreetAddress2:       streetAddress2,
			City:                 city,
			State:                state,
			PostalOrZipCode:      postalOrZipCode,
			GeolocationLatitude:  &latitude,
			GeolocationLongitude: &longitude,
			MaxCapacity:          maxCapacity,
			CenterMode:           centerMode,
			IsNationalCenter:     record[11] == "TRUE",
			CenterURL:            record[12],
		}

		err = db.Create(&newRecord).Error
		if err != nil {
			log.Println("failed to insert record:", err)
		} else {
			log.Println("record inserted successfully:", newRecord)
		}
	}
}
