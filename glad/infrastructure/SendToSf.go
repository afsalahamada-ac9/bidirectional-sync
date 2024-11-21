package Infrastructure

import (
	"bytes"
	"encoding/json"
	utils "glad/common"
	entity "glad/entity"
	"io/ioutil"
	"log"
	"net/http"
)

// this is the data that will be sent from aws
func SendToSF(sendToSf entity.SFData) string {
	sf := utils.GetFromEnv("SEND_DATA_TO_SF")
	jsonData, err := json.Marshal(sendToSf)
	if err != nil {
		log.Println("there is an error in the input file", err)
	}
	req, err := http.NewRequest("POST", sf, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("error creating the request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+utils.GetFromEnv("AUTH_TOKEN"))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("there was an error in the request", err)
	}
	var SfResult string
	parse, err := ioutil.ReadAll(resp.Body)
	log.Println("parse:", string(parse))
	if err != nil {
		log.Println("error decoding the object", err)
	}
	err = json.Unmarshal(parse, &SfResult)
	if err != nil {
		log.Println(err)
	}
	log.Println("data was sent successfully")
	return SfResult
}
