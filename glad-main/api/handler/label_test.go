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

	mock "sudhagar/glad/usecase/label/mock"

	"github.com/codegangsta/negroni"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const (
	labelIDPrimary   entity.ID = 13790492210917030000
	labelIDSecondary entity.ID = 13790492210917030002

	labelNamePrimary   string = "Alpha"
	labelNameSecondary string = "Beta"

	labelColorPrimary   uint32 = 111
	labelColorSecondary uint32 = 222
)

func Test_listLabels(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeLabelHandlers(r, *n, service)
	path, err := r.GetRoute("listLabels").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/labels", path)
	label := &entity.Label{
		ID:       labelIDPrimary,
		TenantID: tenantAlice,
		Name:     labelNamePrimary,
		Color:    labelColorPrimary,
	}
	service.EXPECT().GetCount(label.TenantID).Return(2, nil)
	service.EXPECT().
		GetMulti(label.TenantID, 0, 10).
		Return([]*entity.Label{label, label}, nil)
	ts := httptest.NewServer(listLabels(service))
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	req.Header.Set(httpHeaderTenantID, tenantAlice.String())

	q := req.URL.Query()
	q.Add("index", "0")
	q.Add("limit", "10")
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// presenter.Label is returned by the api (http) server
	var d []presenter.Label
	json.NewDecoder(res.Body).Decode(&d)
	assert.NotNil(t, d)
	assert.Equal(t, 2, len(d))

	assert.Equal(t, label.ID, d[0].ID)
	assert.Equal(t, label.TenantID, d[0].TenantID)
	assert.Equal(t, label.Name, d[0].Name)
	assert.Equal(t, label.Color, d[0].Color)
	assert.Equal(t, tenantAlice.String(), res.Header.Get(httpHeaderTenantID))
}

func Test_createLabel(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeLabelHandlers(r, *n, service)
	path, err := r.GetRoute("createLabel").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/labels", path)

	id := entity.NewID()
	service.EXPECT().
		Create(tenantAlice,
			labelNamePrimary,
			labelColorPrimary).
		Return(id, nil)
	h := createLabel(service)

	ts := httptest.NewServer(h)
	defer ts.Close()

	payload := struct {
		Name  string `json:"name"`
		Color uint32 `json:"color"`
	}{
		Name:  labelNamePrimary,
		Color: labelColorPrimary,
	}
	payloadBytes, err := json.Marshal(payload)
	assert.Nil(t, err)

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost,
		ts.URL+"/v1/labels",
		bytes.NewReader(payloadBytes))
	req.Header.Set(httpHeaderTenantID, tenantAlice.String())
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var label *presenter.Label
	json.NewDecoder(res.Body).Decode(&label)
	assert.Equal(t, id, label.ID)
	assert.Equal(t, payload.Name, label.Name)
	assert.Equal(t, payload.Color, label.Color)
	assert.Equal(t, tenantAlice.String(), res.Header.Get(httpHeaderTenantID))
}

func Test_getLabel(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeLabelHandlers(r, *n, service)
	path, err := r.GetRoute("getLabel").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/labels/{id}", path)

	label := &entity.Label{
		ID:       labelIDPrimary,
		TenantID: tenantAlice,
		Name:     labelNamePrimary,
		Color:    labelColorPrimary,
	}
	service.EXPECT().
		Get(tenantAlice, labelIDPrimary).
		Return(label, nil)
	handler := getLabel(service)

	req, _ := http.NewRequest("GET", "/v1/labels/"+label.ID.String(), nil)
	req.Header.Set(httpHeaderTenantID, tenantAlice.String())

	r.Handle("/v1/labels/{id}", handler).Methods("GET", "OPTIONS")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// presenter.Label is returned by the api (http) server
	var d *presenter.Label
	json.NewDecoder(rr.Body).Decode(&d)
	assert.NotNil(t, d)

	assert.Equal(t, label.ID, d.ID)
	assert.Equal(t, label.Name, d.Name)
	assert.Equal(t, label.Color, d.Color)
	assert.Equal(t, tenantAlice.String(), rr.Header().Get(httpHeaderTenantID))
}

func Test_deleteLabel(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeLabelHandlers(r, *n, service)
	path, err := r.GetRoute("deleteLabel").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/labels/{id}", path)

	id := labelIDPrimary
	service.EXPECT().Delete(tenantAlice, id).Return(nil)
	handler := deleteLabel(service)

	req, _ := http.NewRequest("DELETE", "/v1/labels/"+id.String(), nil)
	req.Header.Set(httpHeaderTenantID, tenantAlice.String())

	r.Handle("/v1/labels/{id}", handler).Methods("DELETE", "OPTIONS")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, tenantAlice.String(), rr.Header().Get(httpHeaderTenantID))
}

func Test_deleteLabelNonExistent(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeLabelHandlers(r, *n, service)
	path, err := r.GetRoute("deleteLabel").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/labels/{id}", path)

	id := labelIDPrimary
	service.EXPECT().Delete(tenantAlice, id).Return(entity.ErrNotFound)
	handler := deleteLabel(service)

	req, _ := http.NewRequest("DELETE", "/v1/labels/"+id.String(), nil)
	req.Header.Set(httpHeaderTenantID, tenantAlice.String())

	r.Handle("/v1/labels/{id}", handler).Methods("DELETE", "OPTIONS")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, tenantAlice.String(), rr.Header().Get(httpHeaderTenantID))
}

// TODO: Test case for updating label
