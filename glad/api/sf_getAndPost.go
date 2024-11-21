package api

import (
	"encoding/json"
	entity "glad/entity"
	infrastructure "glad/infrastructure"
	"io/ioutil"
	"log"
	"net/http"
)

func SendDataToSf(w http.ResponseWriter, r *http.Request) {
	var data entity.SFData
	resp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("there was an error parsing the body")
	}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		log.Println(err)
	}
	response := infrastructure.SendToSF(data)
	json.NewEncoder(w).Encode(data)
	json.NewEncoder(w).Encode(response)
}

func GetDataFromSf(w http.ResponseWriter, r *http.Request) {
	resp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error occurred", err)
		json.NewEncoder(w).Encode(err)
		log.Println(err)
	}
	var result []entity.AWS
	err = json.Unmarshal(resp, &result)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		log.Println(err)
	}
	json.NewEncoder(w).Encode(result)
	log.Println("receive success")
	// todo: postgres save and trigger

}
