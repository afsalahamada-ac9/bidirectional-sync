package main

import (
	"context"
	api "glad/api"
	commonutils "glad/common"
	utils "glad/pkg/utils"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func main() {
	_, err := utils.ConnectDB()
	if err != nil {
		log.Println("error in the database connection", err)
	}
	token, err := utils.CreateTokens()
	if err != nil {
		log.Println("error generating tokens", err)
	}
	err = commonutils.UpdateEnv("AUTH_TOKEN", token)
	if err != nil {
		log.Println(err)
	}
	log.Println("token current:", token)
	go func() {
		ticker := time.NewTicker(60 * time.Minute)
		for {
			<-ticker.C
			log.Println("auth token expired, refreshing the token")
			token, err := utils.CreateTokens()
			if err != nil {
				log.Println("error fetching auth tokens")
			}
			err = commonutils.UpdateEnv("AUTH_TOKEN", token)
			if err != nil {
				log.Println("error updating the auth token")
			}
		}
	}()
	conn, err := pgx.Connect(context.Background(), commonutils.GetFromEnv("CONNECTION_STRING"))
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
	for notification := range channel {
		log.Println("received notification:", notification.Payload)
	}
	log.Println("database connection successful")
	router := mux.NewRouter()
	router.HandleFunc("/sendsf", api.SendDataToSf).Methods("POST")
	router.HandleFunc("/getsf", api.GetDataFromSf).Methods("POST")
	log.Println(http.ListenAndServe(":4000", router))
}
