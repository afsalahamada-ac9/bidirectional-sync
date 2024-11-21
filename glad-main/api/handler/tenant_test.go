/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sudhagar/glad/api/presenter"
	"sudhagar/glad/entity"
	"sudhagar/glad/pkg/common"

	mock "sudhagar/glad/usecase/tenant/mock"

	"github.com/codegangsta/negroni"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const (
	nameAlice    = "alice@wonder.land"
	countryAlice = "testing123"
)

func Test_listTenants(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeTenantHandlers(r, *n, service)
	path, err := r.GetRoute("listTenants").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/tenants", path)
	tenant := &entity.Tenant{
		ID: entity.NewID(),
	}
	service.EXPECT().GetCount().Return(1)
	service.EXPECT().
		ListTenants().
		Return([]*entity.Tenant{tenant}, nil)
	ts := httptest.NewServer(listTenants(service))
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	req.Header.Set(common.HttpHeaderTenantID, tenantAlice.String())
	res, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_createTenant(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeTenantHandlers(r, *n, service)
	path, err := r.GetRoute("createTenant").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/tenants", path)

	id := entity.NewID()
	service.EXPECT().
		CreateTenant(gomock.Any(),
			gomock.Any()).
		Return(id, nil)
	h := createTenant(service)

	ts := httptest.NewServer(h)
	defer ts.Close()

	payload := struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	}{Name: nameAlice,
		Country: countryAlice}
	payloadBytes, err := json.Marshal(payload)
	assert.Nil(t, err)

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost,
		ts.URL+"/v1/tenants",
		bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var tenant *presenter.Tenant
	json.NewDecoder(res.Body).Decode(&tenant)
	// only ID is returned when creating tenant
	assert.Equal(t, id, tenant.ID)
}

func Test_getTenant(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeTenantHandlers(r, *n, service)
	path, err := r.GetRoute("getTenant").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/tenants/{id}", path)
	tenant := &entity.Tenant{
		ID:        entity.NewID(),
		Name:      nameAlice,
		Country:   countryAlice,
		AuthToken: "token123",
	}
	service.EXPECT().
		GetTenant(tenant.ID).
		Return(tenant, nil)
	handler := getTenant(service)
	r.Handle("/v1/tenants/{id}", handler)
	ts := httptest.NewServer(r)
	defer ts.Close()
	res, err := http.Get(ts.URL + "/v1/tenants/" + tenant.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// presenter.Tenant is returned by the api (http) server
	var d *presenter.Tenant
	json.NewDecoder(res.Body).Decode(&d)
	assert.NotNil(t, d)

	assert.Equal(t, tenant.ID, d.ID)
	assert.Equal(t, tenant.Name, d.Name)
	// country is not returned
	assert.Equal(t, tenant.AuthToken, d.AuthToken)
}

func Test_deleteTenant(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeTenantHandlers(r, *n, service)
	path, err := r.GetRoute("deleteTenant").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/tenants/{id}", path)

	id := entity.NewID()
	service.EXPECT().DeleteTenant(id).Return(nil)
	handler := deleteTenant(service)
	req, _ := http.NewRequest("DELETE", "/v1/tenants/"+id.String(), nil)
	r.Handle("/v1/tenants/{id}", handler).Methods("DELETE", "OPTIONS")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func Test_deleteTenantNonExistent(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeTenantHandlers(r, *n, service)
	path, err := r.GetRoute("deleteTenant").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/tenants/{id}", path)

	id := entity.NewID()
	service.EXPECT().DeleteTenant(id).Return(entity.ErrNotFound)
	handler := deleteTenant(service)
	req, _ := http.NewRequest("DELETE", "/v1/tenants/"+id.String(), nil)
	r.Handle("/v1/tenants/{id}", handler).Methods("DELETE", "OPTIONS")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func Test_login(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeTenantHandlers(r, *n, service)
	path, err := r.GetRoute("login").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/login", path)

	tenant := &entity.Tenant{
		ID:        entity.NewID(),
		Name:      nameAlice,
		Country:   countryAlice,
		AuthToken: "token123",
	}
	service.EXPECT().
		Login(tenant.Name, tenant.Country).
		Return(tenant, nil)

	handler := login(service)
	r.Handle("/v1/login", handler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	payload := struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	}{Name: nameAlice,
		Country: countryAlice}
	payloadBytes, err := json.Marshal(payload)
	assert.Nil(t, err)

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost,
		ts.URL+"/v1/login",
		bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// presenter.Tenant is returned by the api (http) server
	var d *presenter.Tenant
	json.NewDecoder(res.Body).Decode(&d)
	assert.NotNil(t, d)

	assert.Equal(t, tenant.ID, d.ID)
	assert.Equal(t, tenant.Name, d.Name)
	// country is not returned
	assert.Equal(t, tenant.AuthToken, d.AuthToken)
}

// TODO: Test case for updating tenant
