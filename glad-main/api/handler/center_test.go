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

	mock "sudhagar/glad/usecase/center/mock"

	"github.com/codegangsta/negroni"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const (
	tenantAlice       entity.ID = 7264348473653242881
	tenantNonExistent entity.ID = 0
	aliceExtID                  = "000aliceExtID"
)

func Test_listCenters(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeCenterHandlers(r, *n, service)
	path, err := r.GetRoute("listCenters").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/centers", path)
	tmpl := &entity.Center{
		ID:       entity.NewID(),
		TenantID: tenantAlice,
		Name:     "default-0",
	}
	service.EXPECT().GetCount(tmpl.TenantID).Return(1)
	service.EXPECT().
		ListCenters(tmpl.TenantID).
		Return([]*entity.Center{tmpl}, nil)
	ts := httptest.NewServer(listCenters(service))
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	req.Header.Set(httpHeaderTenantID, tenantAlice.String())
	res, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_listCenters_NotFound(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	ts := httptest.NewServer(listCenters(service))
	defer ts.Close()
	tenantID := tenantAlice
	service.EXPECT().GetCount(tenantID).Return(0)
	service.EXPECT().
		SearchCenters(tenantID, "non-existent").
		Return(nil, entity.ErrNotFound)

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet,
		ts.URL+"?search=non-existent",
		nil)
	req.Header.Set(httpHeaderTenantID, tenantAlice.String())
	res, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func Test_listCenters_Search(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	tmpl := &entity.Center{
		ID:       entity.NewID(),
		TenantID: tenantAlice,
		Name:     "default-0",
	}
	service.EXPECT().GetCount(tmpl.TenantID).Return(1)
	service.EXPECT().
		SearchCenters(tmpl.TenantID, "default").
		Return([]*entity.Center{tmpl}, nil)
	ts := httptest.NewServer(listCenters(service))
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet,
		ts.URL+"?search=default",
		nil)
	req.Header.Set(httpHeaderTenantID, tenantAlice.String())
	res, err := client.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_createCenter(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeCenterHandlers(r, *n, service)
	path, err := r.GetRoute("createCenter").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/centers", path)

	id := entity.NewID()
	service.EXPECT().
		CreateCenter(gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any()).
		Return(id, nil)
	h := createCenter(service)

	ts := httptest.NewServer(h)
	defer ts.Close()

	payload := struct {
		TenantID entity.ID         `json:"tenant_id"`
		ExtID    string            `json:"extId"`
		Name     string            `json:"name"`
		Mode     entity.CenterMode `json:"mode"`
		Content  string            `json:"content"`
	}{TenantID: tenantAlice,
		ExtID: aliceExtID,
		Name:  "default-0",
		Mode:  (entity.CenterInPerson)}
	payloadBytes, err := json.Marshal(payload)
	assert.Nil(t, err)

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost,
		ts.URL+"/v1/centers",
		bytes.NewReader(payloadBytes))
	req.Header.Set(httpHeaderTenantID, tenantAlice.String())
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var tmpl *presenter.Center
	json.NewDecoder(res.Body).Decode(&tmpl)
	assert.Equal(t, id, tmpl.ID)
	assert.Equal(t, payload.Content, tmpl.Content)
	assert.Equal(t, payload.ExtID, tmpl.ExtID)
	assert.Equal(t, payload.Name, tmpl.Name)
	assert.Equal(t, payload.Mode, tmpl.Mode)
	assert.Equal(t, tenantAlice.String(), res.Header.Get(httpHeaderTenantID))
}

func Test_getCenter(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeCenterHandlers(r, *n, service)
	path, err := r.GetRoute("getCenter").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/centers/{id}", path)
	tmpl := &entity.Center{
		ID:       entity.NewID(),
		TenantID: tenantAlice,
		ExtID:    aliceExtID,
		Name:     "default-0",
		Mode:     entity.CenterInPerson,
	}
	service.EXPECT().
		GetCenter(tmpl.ID).
		Return(tmpl, nil)
	handler := getCenter(service)
	r.Handle("/v1/centers/{id}", handler)
	ts := httptest.NewServer(r)
	defer ts.Close()
	res, err := http.Get(ts.URL + "/v1/centers/" + tmpl.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// presenter.Center is returned by the api (http) server
	var d *presenter.Center
	json.NewDecoder(res.Body).Decode(&d)
	assert.NotNil(t, d)

	assert.Equal(t, tmpl.ID, d.ID)
	assert.Equal(t, tmpl.ExtID, d.ExtID)
	assert.Equal(t, tmpl.Name, d.Name)
	assert.Equal(t, tmpl.Mode, d.Mode)
	assert.Equal(t, tenantAlice.String(), res.Header.Get(httpHeaderTenantID))
}

func Test_deleteCenter(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeCenterHandlers(r, *n, service)
	path, err := r.GetRoute("deleteCenter").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/centers/{id}", path)

	id := entity.NewID()
	service.EXPECT().DeleteCenter(id).Return(nil)
	handler := deleteCenter(service)
	req, _ := http.NewRequest("DELETE", "/v1/centers/"+id.String(), nil)
	r.Handle("/v1/centers/{id}", handler).Methods("DELETE", "OPTIONS")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func Test_deleteCenterNonExistent(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeCenterHandlers(r, *n, service)
	path, err := r.GetRoute("deleteCenter").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/centers/{id}", path)

	id := entity.NewID()
	service.EXPECT().DeleteCenter(id).Return(entity.ErrNotFound)
	handler := deleteCenter(service)
	req, _ := http.NewRequest("DELETE", "/v1/centers/"+id.String(), nil)
	r.Handle("/v1/centers/{id}", handler).Methods("DELETE", "OPTIONS")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// TODO: Test case for updating center
