package utils

import (
	module "glad/Common"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	url := module.GetFromEnv("CONNECTION_STRING")
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Println("there was an error connecting to the database")
		return nil, err
	}
	// err = db.AutoMigrate()
	// if err != nil {
	// 	log.Println("there was an error migrating the database")
	// 	return nil, err
	// }
	return db, nil

}
