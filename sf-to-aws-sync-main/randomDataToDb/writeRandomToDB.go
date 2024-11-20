package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// notes:
// postgres supports LISTEN and NOTIFY commands
// create a trigger in the postgres so that when CUD operations occur, the NOTIFY event is called
// although gorm can be used to interact with the db, it does not support listen and notify events
// so we use pgx library which supports those events

var (
	host     = "localhost"
	port     = 5432
	dbname   = "changeNotifier"
	user     = "postgres"
	password = "50022021"
)

func generateName(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	res := make([]byte, length)
	for i := range res {
		res[i] = letters[rand.Intn(len(letters))]
	}
	return string(res)
}

type Table_name struct {
	ID   int `gorm:"primaryKey"`
	Data string
}

var connection_string = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

func insert() {
	str := generateName(5)
	db, err := gorm.Open(postgres.Open(connection_string), &gorm.Config{})
	if err != nil {
		log.Println("there was an error connecting with the database")
	}
	log.Println("db connected successfully")
	data := Table_name{Data: str}
	res := db.Create(&data)
	if res.Error != nil {
		log.Println("error inserting the record!")
		log.Println("switching to retry mode")
		for i := 0; i < 5; i++ {
			retry_res := db.Create(&data)
			if retry_res.Error == nil {
				log.Println("retry insertion was successful!")
				return
			}
			if i < 4 {
				log.Println("error occurred, trying again")
			}
		}
	}
	log.Println("insertion was successful")
}

func main() {
	// connect to the database
	conn, err := pgx.Connect(context.Background(), connection_string)
	if err != nil {
		log.Println("there was an error when connecting with the database")
	}
	defer conn.Close(context.Background())
	// listen for notifications on events channel -> here the string is a sql command
	_, err = conn.Exec(context.Background(), "LISTEN events")
	if err != nil {
		log.Println("error occurred during listening to the channel")
	}
	log.Println("connection has been setup, now listening on the events channels")

	// the channel to receive the notifications
	notificationChannel := make(chan *pgconn.Notification)

	// goroutine to handle the notifications
	go func() {
		for {
			notification, err := conn.WaitForNotification(context.Background())
			if err != nil {
				log.Println("there was an error listening to the channel")
				continue
			}
			// send the notification to the channel
			notificationChannel <- notification
		}
	}()
	for i := 0; i < 10; i++ {
		insert()
	}
	for notification := range notificationChannel {
		log.Println("notification that was received:", notification, notification.Payload)
	}
}
