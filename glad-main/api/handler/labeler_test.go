/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"sudhagar/glad/entity"
	mock "sudhagar/glad/usecase/labeler/mock"

	"github.com/codegangsta/negroni"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const ()

// TODO: Make tests run in parallel with subtest.Parallel().
// However, bubble up the error so there is an indication that the tests indeed failed.
func Test_SetLabel(t *testing.T) {
	var tests = []struct {
		name    string
		returns error
		want    int
	}{
		{"Happy path", nil, http.StatusOK},
		{"Error", entity.ErrNotFound, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(subtest *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			service := mock.NewMockUseCase(controller)
			r := mux.NewRouter()
			n := negroni.New()
			MakeLabelerHandlers(r, *n, service)

			path, err := r.GetRoute("setLabel").GetPathTemplate()
			assert.Nil(t, err)
			assert.Equal(t, "/v1/contacts/label/{contact_id}/{label_id}", path)

			service.EXPECT().SetLabel(tenantAlice, contactIDAlpha, labelIDPrimary).Return(tt.returns)
			handler := setLabel(service)

			req, _ := http.NewRequest("PUT",
				"/v1/contacts/label/"+contactIDAlpha.String()+"/"+labelIDPrimary.String(),
				nil)
			req.Header.Set(httpHeaderTenantID, tenantAlice.String())

			r.Handle("/v1/contacts/label/{contact_id}/{label_id}", handler).Methods("PUT", "OPTIONS")
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			assert.Equal(subtest, tt.want, rr.Code)
			assert.Equal(subtest, tenantAlice.String(), rr.Header().Get(httpHeaderTenantID))

		})
	}
}

func Test_RemoveLabel(t *testing.T) {
	var tests = []struct {
		name    string
		returns error
		want    int
	}{
		{"Happy path", nil, http.StatusOK},
		{"Error", entity.ErrNotFound, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(subtest *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			service := mock.NewMockUseCase(controller)
			r := mux.NewRouter()
			n := negroni.New()
			MakeLabelerHandlers(r, *n, service)

			path, err := r.GetRoute("removeLabel").GetPathTemplate()
			assert.Nil(t, err)
			assert.Equal(t, "/v1/contacts/label/{contact_id}/{label_id}", path)

			service.EXPECT().RemoveLabel(tenantAlice, contactIDAlpha, labelIDPrimary).Return(tt.returns)
			handler := removeLabel(service)

			req, _ := http.NewRequest("DELETE",
				"/v1/contacts/label/"+contactIDAlpha.String()+"/"+labelIDPrimary.String(),
				nil)
			req.Header.Set(httpHeaderTenantID, tenantAlice.String())

			r.Handle("/v1/contacts/label/{contact_id}/{label_id}", handler).Methods("DELETE", "OPTIONS")
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			assert.Equal(subtest, tt.want, rr.Code)
			assert.Equal(subtest, tenantAlice.String(), rr.Header().Get(httpHeaderTenantID))
		})
	}
}
