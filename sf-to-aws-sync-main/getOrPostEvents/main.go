package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// NOTE: jsonb to string - no.

// TODO:  implement retry mechanism for failing records

// ObjectDescribe represents the "objectDescribe" section of the response
type ObjectDescribe struct {
	Activateable          bool    `json:"activateable"`
	AssociateEntityType   *string `json:"associateEntityType"`
	AssociateParentEntity *string `json:"associateParentEntity"`
	Createable            bool    `json:"createable"`
	Custom                bool    `json:"custom"`
	CustomSetting         bool    `json:"customSetting"`
	DeepCloneable         bool    `json:"deepCloneable"`
	Deletable             bool    `json:"deletable"`
	DeprecatedAndHidden   bool    `json:"deprecatedAndHidden"`
	FeedEnabled           bool    `json:"feedEnabled"`
	HasSubtypes           bool    `json:"hasSubtypes"`
	IsInterface           bool    `json:"isInterface"`
	IsSubtype             bool    `json:"isSubtype"`
	KeyPrefix             string  `json:"keyPrefix"`
	Label                 string  `json:"label"`
	LabelPlural           string  `json:"labelPlural"`
	Layoutable            bool    `json:"layoutable"`
	Mergeable             bool    `json:"mergeable"`
	MruEnabled            bool    `json:"mruEnabled"`
	Name                  string  `json:"name"`
	Queryable             bool    `json:"queryable"`
	Replicateable         bool    `json:"replicateable"`
	Retrieveable          bool    `json:"retrieveable"`
	Searchable            bool    `json:"searchable"`
	Triggerable           bool    `json:"triggerable"`
	Undeletable           bool    `json:"undeletable"`
	Updateable            bool    `json:"updateable"`
	Urls                  struct {
		CompactLayouts  string `json:"compactLayouts"`
		RowTemplate     string `json:"rowTemplate"`
		ApprovalLayouts string `json:"approvalLayouts"`
		Describe        string `json:"describe"`
		QuickActions    string `json:"quickActions"`
		Layouts         string `json:"layouts"`
		SObject         string `json:"sobject"`
	} `json:"urls"`
}

// RecentItem represents an item in the "recentItems" array
type RecentItem struct {
	Attributes struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"attributes"`
	ID   string `json:"Id"`
	Name string `json:"Name"`
}

// SalesforceResponse represents the entire response from Salesforce
type SalesforceResponse struct {
	ObjectDescribe ObjectDescribe `json:"objectDescribe"`
	RecentItems    []RecentItem   `json:"recentItems"`
}

// Event represents the structure of the event data from Salesforce
type Event struct {
	MaxAttendees              int    `json:"Max_Attendees__c"`
	RegistrationStartDateTime string `json:"Registration_Start_Date_Time__c"`
	RegistrationEndDateTime   string `json:"Registration_End_Date_Time__c"`
	Location                  string `json:"Location__c"`
	EventEndDate              string `json:"Event_End_Date__c"`
	EventStartDate            string `json:"Event_Start_Date__c"`
	Notes                     string `json:"Notes__c"`
	Status                    string `json:"Status__c"`
	EventStartTime            string `json:"Event_Start_Time__c"`
	EventEndTime              string `json:"Event_End_Time__c"`
	Timezone                  string `json:"Timezone__c"`
	ContactPerson             string `json:"Contact_Person__c"`
	Organizer                 string `json:"Organizer__c"`
	EventStartDateTimeGMT     string `json:"Event_Start_Date_Time_GMT__c"`
	EventEndDateTimeGMT       string `json:"Event_End_Date_Time_GMT__c"`
}

type TokenResponse struct {
	Accesstoken string `json:"access_token"`
	InstanceURL string `json:"instance_url"`
	ExpiresIn   string `json:"expires_in"`
}

var AccessToken_Global string

// postEventRecord sends the event data to the Salesforce endpoint
func postEventRecord(event []Event) ([]byte, error) {
	url := "https://aol-dev--awspoc.sandbox.my.salesforce.com/services/apexrest/handleAolEvent"
	jsonData, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+AccessToken_Global)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("failed to post event: %s", body)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// TODO: this will handle the post request to post an event record
func handlePostEventRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var events []Event
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&events)
	if err != nil {
		log.Println("an error occurred when decoding the event", err)
	}
	result, err := postEventRecord(events)
	if err != nil {
		log.Println("there was an issue posting the event record", err)
	}
	// json.NewEncoder(w).Encode(events)
	var apiResponse string
	err = json.Unmarshal(result, &apiResponse)
	if err != nil {
		log.Println("error decoding the json", err)
		return
	}
	json.NewEncoder(w).Encode(apiResponse)
}

// handleGetEvent handles the GET request for events
func handleGetEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// bearerToken := "00DVZ000001OUfV!AQEAQDInDzxaaWo7AZglE2z1OdzeFrA6pzaKj8SjIWEnaGkU2trqYornUaaVaYYlHWQUGX8NZVMxPfK4lpiQQMG9_.mNdCA8"

	// Create a new request with the Bearer token
	req, err := http.NewRequest("GET", "https://aol-dev--awspoc.sandbox.my.salesforce.com/services/data/v55.0/sobjects/Event__c", nil)
	if err != nil {
		log.Println("error creating request:", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+AccessToken_Global)

	// Make the HTTP GET request
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error fetching the response:", err)
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close() // Ensure the response body is closed

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		log.Printf("Received non-200 response: %d\n", response.StatusCode)
		body, _ := ioutil.ReadAll(response.Body) // Read the body for logging
		log.Printf("Response body: %s\n", body)
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("error reading the response body:", err)
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// Log the response body for debugging
	log.Printf("Response body: %s\n", body)

	// Dynamically parse the JSON response
	var parsedResponse SalesforceResponse
	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	// Respond with the parsed data
	json.NewEncoder(w).Encode(parsedResponse)
}

func getAccessTokens() (string, error) {
	conn_url := "https://aol-dev--awspoc.sandbox.my.salesforce.com/services/oauth2/token"
	clientId := getValue("CLIENT_ID")
	clientSecret := getValue("CLIENT_SECRET")
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
	req, err := http.NewRequest("POST", conn_url, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("there was an error processing the request")
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		return "", fmt.Errorf("failed to get the access tokens", body)
	}

	var tokenResponse TokenResponse
	err = json.NewDecoder(response.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}
	accessToken := tokenResponse.Accesstoken
	// tokenExpiresAt := time.Now().Add(time.Duration(tokenResponse.ExpiresIn)*time.Second)
	return accessToken, nil
}

func getValue(key string) string {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("error reading the .env file")
	}
	value, ok := viper.Get(key).(string)
	if !ok {
		log.Println("invalid configuration")
	}
	return value

}

func main() {
	// AccessToken_Global, _ = getAccessTokens()
	// log.Println(AccessToken_Global)
	// FIXME: make this more secure, and for now this is hardcoded in sf db, we can configure this later as required. uncomment above two lines to get it dynamically.
	go func() {
		AccessToken_Global = getValue("ACCESS_TOKEN")
		ticker := time.NewTicker(60 * time.Minute)
		defer ticker.Stop()
		for {
			<-ticker.C
			AccessToken_Global, _ = getAccessTokens()
			// log.Println("tokens were updated", AccessToken_Global)
			log.Println("tokens were updated")
		}
	}()
	router := mux.NewRouter()
	router.HandleFunc("/get", handleGetEvent).Methods("GET")
	router.HandleFunc("/postEvent", handlePostEventRecord).Methods("POST")
	log.Println(http.ListenAndServe(":4000", router))
}
