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

	mock "sudhagar/glad/usecase/contact/mock"

	"github.com/codegangsta/negroni"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const (
	contactIDAlpha entity.ID = 13790492210917050000
	contactIDBeta  entity.ID = 13790492210917050002

	contactHandleAlpha = "12345555001"
	contactHandleBeta  = "12345555002"

	contactNameAlpha = "Alpha Alice"
	contactNameBeta  = "Beta Bob"
)

// TODO: Add test cases for:
// a) Default index and limit values
// b) Invalid limit value
// c) Custom index and limit values (DONE)
func Test_listContacts(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeContactHandlers(r, *n, service)
	path, err := r.GetRoute("listContacts").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/contacts", path)
	tmpl := &entity.Contact{
		ID:        contactIDAlpha,
		TenantID:  tenantAlice,
		AccountID: accountIDPrimary,
		Handle:    contactHandleAlpha,
		Name:      contactNameAlpha,
	}
	service.EXPECT().GetCount(tmpl.TenantID).Return(2, nil)
	service.EXPECT().
		GetMulti(tmpl.TenantID, 0, 10).
		Return([]*entity.Contact{tmpl, tmpl}, nil)
	ts := httptest.NewServer(listContacts(service))
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

	// presenter.Contact is returned by the api (http) server
	var d []presenter.Contact
	json.NewDecoder(res.Body).Decode(&d)
	assert.NotNil(t, d)
	assert.Equal(t, 2, len(d))

	assert.Equal(t, tmpl.ID, d[0].ID)
	assert.Equal(t, tmpl.TenantID, d[0].TenantID)
	assert.Equal(t, tmpl.AccountID, d[0].AccountID)
	assert.Equal(t, tmpl.Handle, d[0].Handle)
	assert.Equal(t, tmpl.Name, d[0].Name)
	assert.Equal(t, tenantAlice.String(), res.Header.Get(httpHeaderTenantID))
}
