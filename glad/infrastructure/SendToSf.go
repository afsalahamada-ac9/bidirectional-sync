package Infrastructure

import (
	"bytes"
	"encoding/json"
	utils "glad/common"
	entity "glad/entity"
	"log"
	"net/http"
)

// this is the data that will be sent from aws
func SendToSF(sendToSf entity.SFData) {
	sf := utils.GetFromEnv("SEND_DATA_TO_SF")
	jsonData, err := json.Marshal(sendToSf)
	if err != nil {
		log.Println("there is an error in the input file", err)
	}
	req, err := http.NewRequest("POST", sf, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("error creating the request")
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("there was an error in the request", err)
	}
	var DataToBeSent entity.SFData
	err = json.NewDecoder(resp.Body).Decode(&DataToBeSent)
	if err != nil {
		log.Println("error decoding the object")
	}
	log.Println("data was sent successfully")
}
