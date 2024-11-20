package Infrastructure

import (
	"encoding/json"
	common "glad/Common"
	entity "glad/entity"
	"log"
	"net/http"
)

// note: this represents data that is sent from salesforce
func FetchSfData() entity.AWS {
	apiUrl := common.GetFromEnv("GET_SF_DATA")
	body := "Bearer " + common.GetFromEnv("AUTH_TOKEN")
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
	var sf_response entity.AWS
	err = json.NewDecoder(resp.Body).Decode(&sf_response)
	defer resp.Body.Close()
	if err != nil {
		log.Println("there was an error parsing the response")
	}
	return sf_response
}
