package Infrastructure

import (
	"encoding/json"
	utils "glad/common"
	entity "glad/entity"
	"io/ioutil"
	"log"
	"net/http"
)

type SfResponse struct {
	Sf_response []entity.AWS `json:"response"`
}

// note: this represents data that is sent from salesforce
func FetchSfData() []entity.AWS {
	apiUrl := utils.GetFromEnv("GET_SF_DATA")
	token := utils.GetFromEnv("AUTH_TOKEN")
	log.Println("token:", token)
	body := "Bearer " + token
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Println("there was an error creating the request")
	}
	req.Header.Set("Authorization", body)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("there was an error processing the request", err)
	}
	response, err := ioutil.ReadAll(resp.Body)
	var result SfResponse
	err = json.Unmarshal(response, &result)
	defer resp.Body.Close()
	if err != nil {
		log.Println("there was an error parsing the response", err)
	}
	return result.Sf_response
}
