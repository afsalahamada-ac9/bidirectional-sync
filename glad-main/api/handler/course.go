/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"sudhagar/glad/usecase/course"

	"sudhagar/glad/api/presenter"

	"sudhagar/glad/entity"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

const ()

func listCourses(service course.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading courses"
		var data []*entity.Course
		var err error
		tenant := r.Header.Get(httpHeaderTenantID)
		search := r.URL.Query().Get("search")

		tenantID, err := entity.StringToID(tenant)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unable to parse tenant id"))
			return
		}

		switch {
		case search == "":
			data, err = service.ListCourses(tenantID)
		default:
			// TODO: search need to be reworked; need to add a count
			// for search; also need to see how the caller generates
			// the search query request
			data, err = service.SearchCourses(tenantID, search)
		}
		w.Header().Set("Content-Type", "application/json")
		if err != nil && err != entity.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage + ":" + err.Error()))
			return
		}

		total := service.GetCount(tenantID)
		w.Header().Set(httpHeaderTotalCount, strconv.Itoa(total))

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}
		var toJ []*presenter.Course
		for _, d := range data {
			toJ = append(toJ, &presenter.Course{
				ID:       d.ID,
				TenantID: d.TenantID,
				ExtID:    d.ExtID,
				Name:     d.Name,
				CType:    d.CType,
			})
		}
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.Header().Set(httpHeaderTenantID, tenant)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to encode course"))
		}
	})
}

func createCourse(service course.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error adding course"
		var input struct {
			ExtID        string                `json:"extId"`
			CenterID     entity.ID             `json:"centerId"`
			Name         string                `json:"name"`
			Notes        string                `json:"notes"`
			Timezone     string                `json:"timezone"`
			Location     entity.CourseLocation `json:"location"`
			Status       entity.CourseStatus   `json:"status"`
			CType        entity.CourseType     `json:"ctype"`
			MaxAttendees int32                 `json:"maxAttendees"`
			// TODO: Organizer, Teacher, Contact, Notify, etc
		}

		tenant := r.Header.Get(httpHeaderTenantID)
		tenantID, err := entity.StringToID(tenant)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing tenant ID"))
			return
		}

		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unable to decode the data. " + err.Error()))
			return
		}

		id, err := service.CreateCourse(
			tenantID,
			input.ExtID,
			input.CenterID,
			input.Name,
			input.Notes,
			input.Timezone,
			input.Location,
			input.Status,
			input.CType,
			input.MaxAttendees,
			// numAttendees, isAutoApprove
			0, false)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage + ":" + err.Error()))
			return
		}
		toJ := &presenter.Course{
			ID:       id,
			ExtID:    input.ExtID,
			Name:     input.Name,
			TenantID: tenantID,
		}

		w.Header().Set(httpHeaderTenantID, tenant)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

func getCourse(service course.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading course"
		vars := mux.Vars(r)
		id, err := entity.StringToID(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		data, err := service.GetCourse(id)
		if err != nil && err != entity.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage + ":" + err.Error()))
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Empty data returned"))
			return
		}

		toJ := &presenter.Course{
			ID:    data.ID,
			ExtID: data.ExtID,
			Name:  data.Name,
			// TODO: More fields to be added
		}

		w.Header().Set(httpHeaderTenantID, data.TenantID.String())
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to encode course"))
		}
	})
}

func deleteCourse(service course.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error removing course"
		vars := mux.Vars(r)
		id, err := entity.StringToID(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		err = service.DeleteCourse(id)
		switch err {
		case nil:
			w.WriteHeader(http.StatusOK)
			return
		case entity.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Course doesn't exist"))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

func updateCourse(service course.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error updating course"

		vars := mux.Vars(r)
		id, err := entity.StringToID(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		var input entity.Course
		tenant := r.Header.Get(httpHeaderTenantID)
		tenantID, err := entity.StringToID(tenant)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing tenant ID"))
			return
		}

		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unable to decode the data. " + err.Error()))
			return
		}

		input.ID = id
		input.TenantID = tenantID
		err = service.UpdateCourse(&input)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage + ":" + err.Error()))
			return
		}

		toJ := &presenter.Course{
			ID:       input.ID,
			TenantID: tenantID,
			ExtID:    input.ExtID,
			Name:     input.Name,
		}

		w.Header().Set(httpHeaderTenantID, tenant)
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

// MakeCourseHandlers make url handlers
func MakeCourseHandlers(r *mux.Router, n negroni.Negroni, service course.UseCase) {
	r.Handle("/v1/courses", n.With(
		negroni.Wrap(listCourses(service)),
	)).Methods("GET", "OPTIONS").Name("listCourses")

	r.Handle("/v1/courses", n.With(
		negroni.Wrap(createCourse(service)),
	)).Methods("POST", "OPTIONS").Name("createCourse")

	r.Handle("/v1/courses/{id}", n.With(
		negroni.Wrap(getCourse(service)),
	)).Methods("GET", "OPTIONS").Name("getCourse")

	r.Handle("/v1/courses/{id}", n.With(
		negroni.Wrap(deleteCourse(service)),
	)).Methods("DELETE", "OPTIONS").Name("deleteCourse")

	r.Handle("/v1/courses/{id}", n.With(
		negroni.Wrap(updateCourse(service)),
	)).Methods("PUT", "OPTIONS").Name("updateCourse")
}
