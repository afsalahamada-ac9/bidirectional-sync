package utils

import (
	"bytes"
	"encoding/json"
	common "glad/common"
	token "glad/entity"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func CreateTokens() (string, error) {
	log.Println("token updated")
	client_id := common.GetFromEnv("CLIENT_ID")
	client_secret := common.GetFromEnv("CLIENT_SECRET")
	apiUrl := common.GetFromEnv("AUTH_TOKEN_API")
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", client_id)
	data.Set("client_secret", client_secret)
	log.Println(data.Encode())
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
	parser, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("error in parse", err)
	}
	log.Println(string(parser))
	err = json.Unmarshal(parser, &result)
	if err != nil {
		return "", err
	}
	log.Println("create token:", result.AuthToken)
	return result.AuthToken, nil
}
