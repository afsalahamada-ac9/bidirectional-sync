package main

import (
	utils "glad/pkg/utils"
	"log"
)

func main() {
	_, err := utils.ConnectDB()
	if err != nil {
		log.Println("error in the database connection", err)
	}
	// go func() {
	// 	ticker := time.NewTicker(60 * time.Minute)
	// 	for {
	// 		<-ticker.C
	// 		log.Println("auth token expired, refreshing the token")
	// 		token, err := utils.CreateTokens()
	// 		if err != nil {
	// 			log.Println("error fetching auth tokens")
	// 		}
	// 		err = common.WriteToEnv("AUTH_TOKEN", token)
	// 		if err != nil {
	// 			log.Println("error updating the auth token")
	// 		}
	// 	}
	// }()
	log.Println("database connection successful")
}
