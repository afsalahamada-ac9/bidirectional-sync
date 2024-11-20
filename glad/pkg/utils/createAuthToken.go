package utils

import (
	"bytes"
	"encoding/json"
	common "glad/Common"
	token "glad/entity"
	"log"
	"net/http"
	"net/url"
)

func CreateTokens() (string, error) {
	client_id := common.GetFromEnv("CLIENT_ID")
	client_secret := common.GetFromEnv("CLIENT_SECRET")
	apiUrl := common.GetFromEnv("AUTH_TOKEN_API")
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", client_id)
	data.Set("client_secret", client_secret)
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Println("error creating the tokens")
		return "", err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error in request", err)
		return "", err
	}
	var result token.Token
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result.AuthToken, nil
}
