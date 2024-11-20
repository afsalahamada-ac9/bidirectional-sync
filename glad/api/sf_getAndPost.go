package api

import (
	"encoding/json"
	entity "glad/entity"
	infrastructure "glad/infrastructure"
	"log"
	"net/http"
)

func sendDataToSf(w http.ResponseWriter, r *http.Request) {
	var data entity.SFData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("there was an error parsing the body")
	}
	infrastructure.SendToSF(data)
}

func getDataFromSf(w http.ResponseWriter, r *http.Request) {
	resp := infrastructure.FetchSfData()
	json.NewEncoder(w).Encode(resp)
}
