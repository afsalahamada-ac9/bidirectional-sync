package main

import (
	api "glad/api"
	commonutils "glad/common"
	utils "glad/pkg/utils"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

	log.Println("database connection successful")
	router := mux.NewRouter()
	router.HandleFunc("/sendsf", api.SendDataToSf)
	router.HandleFunc("/getsf", api.GetDataFromSf)
	log.Println(http.ListenAndServe(":4000", router))
}
