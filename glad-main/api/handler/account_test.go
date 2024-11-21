/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sudhagar/glad/api/presenter"
	"sudhagar/glad/entity"
	"sudhagar/glad/pkg/common"

	mock "sudhagar/glad/usecase/account/mock"

	"github.com/codegangsta/negroni"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const (
	accountIDPrimary   entity.ID = 13790492210917010000
	accountIDSecondary entity.ID = 13790492210917010002

	accountUsernamePrimary   string = "12345550001"
	accountUsernameSecondary string = "12345550002"
)

func Test_listAccounts(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeAccountHandlers(r, *n, service)
	path, err := r.GetRoute("listAccounts").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/accounts", path)
	tmpl := &entity.Account{
		ID:       accountIDPrimary,
		ExtID:    aliceExtID,
		Username: accountUsernamePrimary,
		Type:     entity.AccountTeacher,
	}
	service.EXPECT().GetCount(tenantAlice).Return(1)
	service.EXPECT().
		ListAccounts(tenantAlice).
		Return([]*entity.Account{tmpl}, nil)
	ts := httptest.NewServer(listAccounts(service))
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	req.Header.Set(common.HttpHeaderTenantID, tenantAlice.String())
	res, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

// func Test_createAccount(t *testing.T) {
// 	controller := gomock.NewController(t)
// 	defer controller.Finish()
// 	service := mock.NewMockUseCase(controller)
// 	r := mux.NewRouter()
// 	n := negroni.New()
// 	MakeAccountHandlers(r, *n, service)
// 	path, err := r.GetRoute("createAccount").GetPathTemplate()
// 	assert.Nil(t, err)
// 	assert.Equal(t, "/v1/accounts", path)

// 	id := entity.NewID()
// 	service.EXPECT().
// 		CreateAccount(gomock.Any(),
// 			gomock.Any(),
// 			gomock.Any(),
// 			gomock.Any()).
// 		Return(id, nil)
// 	h := createAccount(service)

// 	ts := httptest.NewServer(h)
// 	defer ts.Close()

// 	payload := struct {
// 		TenantID entity.ID          `json:"tenant_id"`
// 		Username     string             `json:"Username"`
// 		Type     entity.AccountType `json:"type"`
// 		Content  string             `json:"content"`
// 	}{TenantID: tenantAlice,
// 		Username:    "default-0",
// 		Type:    (entity.AccountText),
// 		Content: "This is a default message"}
// 	payloadBytes, err := json.Marshal(payload)
// 	assert.Nil(t, err)

// 	client := &http.Client{}
// 	req, _ := http.NewRequest(http.MethodPost,
// 		ts.URL+"/v1/accounts",
// 		bytes.NewReader(payloadBytes))
// 	req.Header.Set(common.HttpHeaderTenantID, tenantAlice)
// 	req.Header.Set("Content-Type", "application/json")
// 	res, err := client.Do(req)

// 	assert.Nil(t, err)
// 	assert.Equal(t, http.StatusCreated, res.StatusCode)

// 	var tmpl *presenter.Account
// 	json.NewDecoder(res.Body).Decode(&tmpl)
// 	assert.Equal(t, id, tmpl.ID)
// 	assert.Equal(t, payload.Content, tmpl.Content)
// 	assert.Equal(t, payload.Username, tmpl.Username)
// 	assert.Equal(t, payload.Type, tmpl.Type)
// 	assert.Equal(t, tenantAlice, res.Header.Get(common.HttpHeaderTenantID))
// }

func Test_getAccount(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeAccountHandlers(r, *n, service)
	path, err := r.GetRoute("getAccount").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/accounts/{username}", path)
	tmpl := &entity.Account{
		ID:       accountIDPrimary,
		ExtID:    aliceExtID,
		Username: accountUsernamePrimary,
		Type:     entity.AccountTeacher,
	}
	service.EXPECT().
		GetAccount(tmpl.Username).
		Return(tmpl, nil)
	handler := getAccount(service)
	r.Handle("/v1/accounts/{username}", handler)
	ts := httptest.NewServer(r)
	defer ts.Close()
	res, err := http.Get(ts.URL + "/v1/accounts/" + tmpl.Username)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// presenter.Account is returned by the api (http) server
	var d *presenter.Account
	json.NewDecoder(res.Body).Decode(&d)
	assert.NotNil(t, d)

	assert.Equal(t, tmpl.ID, d.ID)
	assert.Equal(t, tmpl.Username, d.Username)
	assert.Equal(t, tmpl.Type, d.Type)
	// assert.Equal(t, tenantAlice.String(), res.Header.Get(common.HttpHeaderTenantID))
}

func Test_deleteAccount(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeAccountHandlers(r, *n, service)
	path, err := r.GetRoute("deleteAccount").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/accounts/{username}", path)

	username := accountUsernamePrimary
	service.EXPECT().DeleteAccount(username).Return(nil)
	handler := deleteAccount(service)
	req, _ := http.NewRequest("DELETE", "/v1/accounts/"+username, nil)
	r.Handle("/v1/accounts/{username}", handler).Methods("DELETE", "OPTIONS")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func Test_deleteAccountNonExistent(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeAccountHandlers(r, *n, service)
	path, err := r.GetRoute("deleteAccount").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/accounts/{username}", path)

	username := accountUsernamePrimary
	service.EXPECT().DeleteAccount(username).Return(entity.ErrNotFound)
	handler := deleteAccount(service)
	req, _ := http.NewRequest("DELETE", "/v1/accounts/"+username, nil)
	r.Handle("/v1/accounts/{username}", handler).Methods("DELETE", "OPTIONS")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// TODO: Test case for updating account
