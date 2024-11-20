package utils

import (
	commonFiles "glad/common"
	entity "glad/entity"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	// url := "host=localhost port=5432 user=postgres password=1234 dbname=postgres sslmode=disable"
	url := commonFiles.GetFromEnv("CONNECTION_STRING")
	log.Println(url, "is url")
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Println("there was an error connecting to the database", err)
		return nil, err
	}
	err = db.AutoMigrate(&entity.EventValue{})
	if err != nil {
		log.Println("there was an error migrating the database", err)
		return nil, err
	}
	log.Println("migration success")
	return db, nil

}
