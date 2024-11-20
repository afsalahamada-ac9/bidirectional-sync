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

	mock "sudhagar/glad/usecase/course/mock"

	"github.com/codegangsta/negroni"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const (
	aliceCenterID = 13790493495087075501
)

func Test_listCourses(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeCourseHandlers(r, *n, service)
	path, err := r.GetRoute("listCourses").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/courses", path)
	tmpl := &entity.Course{
		ID:       entity.NewID(),
		TenantID: tenantAlice,
		Name:     "default-0",
	}
	service.EXPECT().GetCount(tmpl.TenantID).Return(1)
	service.EXPECT().
		ListCourses(tmpl.TenantID).
		Return([]*entity.Course{tmpl}, nil)
	ts := httptest.NewServer(listCourses(service))
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	req.Header.Set(httpHeaderTenantID, tenantAlice.String())
	res, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_listCourses_NotFound(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	ts := httptest.NewServer(listCourses(service))
	defer ts.Close()
	tenantID := tenantAlice
	service.EXPECT().GetCount(tenantID).Return(0)
	service.EXPECT().
		SearchCourses(tenantID, "non-existent").
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

func Test_listCourses_Search(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	tmpl := &entity.Course{
		ID:       entity.NewID(),
		TenantID: tenantAlice,
		Name:     "default-0",
	}
	service.EXPECT().GetCount(tmpl.TenantID).Return(1)
	service.EXPECT().
		SearchCourses(tmpl.TenantID, "default").
		Return([]*entity.Course{tmpl}, nil)
	ts := httptest.NewServer(listCourses(service))
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

func Test_createCourse(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeCourseHandlers(r, *n, service)
	path, err := r.GetRoute("createCourse").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/courses", path)

	id := entity.NewID()
	service.EXPECT().
		CreateCourse(gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any()).
		Return(id, nil)
	h := createCourse(service)

	ts := httptest.NewServer(h)
	defer ts.Close()

	payload := struct {
		TenantID entity.ID         `json:"tenant_id"`
		ExtID    string            `json:"extId"`
		Name     string            `json:"name"`
		CType    entity.CourseType `json:"ctype"`
		// CenterID entity.ID         `json:"center_id"`
	}{TenantID: tenantAlice,
		ExtID: aliceExtID,
		Name:  "default-0",
		CType: (entity.CourseInPerson),
		// CenterID: aliceCenterID,
	}
	payloadBytes, err := json.Marshal(payload)
	assert.Nil(t, err)

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost,
		ts.URL+"/v1/courses",
		bytes.NewReader(payloadBytes))
	req.Header.Set(httpHeaderTenantID, tenantAlice.String())
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var tmpl *presenter.Course
	json.NewDecoder(res.Body).Decode(&tmpl)
	assert.Equal(t, id, tmpl.ID)
	assert.Equal(t, payload.ExtID, tmpl.ExtID)
	assert.Equal(t, payload.Name, tmpl.Name)
	assert.Equal(t, payload.CType, tmpl.CType)
	assert.Equal(t, tenantAlice.String(), res.Header.Get(httpHeaderTenantID))
}

func Test_getCourse(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeCourseHandlers(r, *n, service)
	path, err := r.GetRoute("getCourse").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/courses/{id}", path)
	tmpl := &entity.Course{
		ID:       entity.NewID(),
		TenantID: tenantAlice,
		ExtID:    aliceExtID,
		Name:     "default-0",
		CType:    entity.CourseInPerson,
	}
	service.EXPECT().
		GetCourse(tmpl.ID).
		Return(tmpl, nil)
	handler := getCourse(service)
	r.Handle("/v1/courses/{id}", handler)
	ts := httptest.NewServer(r)
	defer ts.Close()
	res, err := http.Get(ts.URL + "/v1/courses/" + tmpl.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// presenter.Course is returned by the api (http) server
	var d *presenter.Course
	json.NewDecoder(res.Body).Decode(&d)
	assert.NotNil(t, d)

	assert.Equal(t, tmpl.ID, d.ID)
	assert.Equal(t, tmpl.ExtID, d.ExtID)
	assert.Equal(t, tmpl.Name, d.Name)
	assert.Equal(t, tmpl.CType, d.CType)
	assert.Equal(t, tenantAlice.String(), res.Header.Get(httpHeaderTenantID))
}

func Test_deleteCourse(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeCourseHandlers(r, *n, service)
	path, err := r.GetRoute("deleteCourse").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/courses/{id}", path)

	id := entity.NewID()
	service.EXPECT().DeleteCourse(id).Return(nil)
	handler := deleteCourse(service)
	req, _ := http.NewRequest("DELETE", "/v1/courses/"+id.String(), nil)
	r.Handle("/v1/courses/{id}", handler).Methods("DELETE", "OPTIONS")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func Test_deleteCourseNonExistent(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeCourseHandlers(r, *n, service)
	path, err := r.GetRoute("deleteCourse").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/courses/{id}", path)

	id := entity.NewID()
	service.EXPECT().DeleteCourse(id).Return(entity.ErrNotFound)
	handler := deleteCourse(service)
	req, _ := http.NewRequest("DELETE", "/v1/courses/"+id.String(), nil)
	r.Handle("/v1/courses/{id}", handler).Methods("DELETE", "OPTIONS")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// TODO: Test case for updating course
