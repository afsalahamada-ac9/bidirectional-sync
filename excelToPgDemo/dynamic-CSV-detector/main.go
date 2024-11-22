package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DynamicRow represents a dynamic data row
type DynamicRow map[string]interface{}

func main() {
	// PostgreSQL connection details
	dsn := "host=localhost user=postgres password=1234 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Directory to monitor
	watchDir := "./test"

	// Create the directory if it doesn't exist
	if _, err := os.Stat(watchDir); os.IsNotExist(err) {
		err := os.Mkdir(watchDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
	}

	// Create a new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Add the directory to the watcher
	err = watcher.Add(watchDir)
	if err != nil {
		log.Fatalf("Failed to watch directory: %v", err)
	}

	log.Printf("Watching directory: %s", watchDir)

	// Start watching for file events
	for {
		select {
		case event := <-watcher.Events:
			// Handle file creation events
			if event.Op&fsnotify.Create == fsnotify.Create {
				if strings.HasSuffix(event.Name, ".csv") {
					log.Printf("Detected new CSV file: %s", event.Name)
					err := processCSV(event.Name, db)
					if err != nil {
						log.Printf("Error processing CSV file %s: %v", event.Name, err)
					} else {
						log.Printf("Successfully processed and inserted data from: %s", filepath.Base(event.Name))
					}
				}
			}
		case err := <-watcher.Errors:
			log.Printf("Watcher error: %v", err)
		}
	}
}

// processCSV reads the CSV file, creates the table dynamically, and inserts the data
func processCSV(filePath string, db *gorm.DB) error {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Read the CSV content
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV content: %v", err)
	}

	// Check if the CSV is valid (has at least headers and one data row)
	if len(records) < 2 {
		return fmt.Errorf("CSV file is empty or does not have enough data")
	}

	// Get the table name and column headers
	tableName := "dynamic_table" // Use a fixed table name or derive it dynamically
	headers := records[0]

	// Create the dynamic table
	err = createTable(db, tableName, headers)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	// Insert the rows into the database
	err = insertData(db, tableName, headers, records[1:])
	if err != nil {
		return fmt.Errorf("failed to insert data: %v", err)
	}

	return nil
}

// createTable creates or updates the table structure based on the CSV headers
func createTable(db *gorm.DB, tableName string, headers []string) error {
	// Build the SQL for dynamic table creation
	columnDefinitions := []string{}
	for _, header := range headers {
		columnName := strings.ToLower(strings.ReplaceAll(header, " ", "_"))
		columnDefinitions = append(columnDefinitions, fmt.Sprintf(`"%s" TEXT`, columnName)) // Assuming TEXT for simplicity
	}

	createTableQuery := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (%s);`,
		tableName, strings.Join(columnDefinitions, ", "),
	)

	// Execute the table creation query
	return db.Exec(createTableQuery).Error
}

// insertData inserts the data into the database dynamically
func insertData(db *gorm.DB, tableName string, headers []string, rows [][]string) error {
	for _, row := range rows {
		if len(row) != len(headers) {
			log.Printf("Skipping row due to mismatch in column count: %v", row)
			continue
		}

		// Build a map of column-value pairs for insertion
		data := DynamicRow{}
		for i, header := range headers {
			columnName := strings.ToLower(strings.ReplaceAll(header, " ", "_"))
			data[columnName] = row[i]
		}

		// Build the INSERT query
		columns := []string{}
		values := []string{}
		for column, value := range data {
			columns = append(columns, fmt.Sprintf(`"%s"`, column))
			values = append(values, fmt.Sprintf("'%s'", value))
		}

		insertQuery := fmt.Sprintf(
			`INSERT INTO %s (%s) VALUES (%s);`,
			tableName, strings.Join(columns, ", "), strings.Join(values, ", "),
		)

		// Execute the INSERT query
		err := db.Exec(insertQuery).Error
		if err != nil {
			return fmt.Errorf("failed to insert row: %v", err)
		}
	}

	return nil
}
