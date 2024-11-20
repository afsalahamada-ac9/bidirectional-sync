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

	"sudhagar/glad/usecase/label"

	"sudhagar/glad/api/presenter"

	"sudhagar/glad/entity"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

const (
	apiLabelListDefaultLimit int = 32
	apiLabelListMaxLimit     int = 128
)

func listLabels(service label.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading labels"
		var data []*entity.Label
		var err error
		tenant := r.Header.Get(httpHeaderTenantID)

		tenantID, err := entity.StringToID(tenant)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unable to parse tenant id"))
			return
		}

		// get page (index) and page_size (limit) from query parameters
		queryParams := r.URL.Query()
		qpIndex := queryParams.Get(apiQueryParamKeyIndex)
		qpLimit := queryParams.Get(apiQueryParamKeyLimit)

		index := 0
		limit := apiLabelListDefaultLimit

		if qpIndex != "" {
			if index, err = strconv.Atoi(qpIndex); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Unable to parse index parameter"))
				return
			}
		}

		if qpLimit != "" {
			limit, err = strconv.Atoi(qpLimit)
			if err != nil || limit > apiLabelListMaxLimit {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Unable to parse limit parameter or invalid limit value"))
				return
			}
		}

		data, err = service.GetMulti(tenantID, index, limit)
		w.Header().Set("Content-Type", "application/json")
		if err != nil && err != entity.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage + ":" + err.Error()))
			return
		}

		total, err := service.GetCount(tenantID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage + ":" + err.Error()))
			return
		}
		w.Header().Set(httpHeaderTotalCount, strconv.Itoa(total))
		w.Header().Set(httpHeaderTenantID, tenant)

		var toJ []*presenter.Label
		for _, d := range data {
			toJ = append(toJ, &presenter.Label{
				ID:       d.ID,
				TenantID: d.TenantID,
				Name:     d.Name,
				Color:    d.Color,
			})
		}
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to encode label"))
		}
	})
}

func createLabel(service label.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error adding label"
		var input struct {
			Name  string `json:"name"`
			Color uint32 `json:"color"`
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

		id, err := service.Create(
			tenantID,
			input.Name,
			input.Color)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage + ":" + err.Error()))
			return
		}
		toJ := &presenter.Label{
			ID:       id,
			Name:     input.Name,
			Color:    input.Color,
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

func getLabel(service label.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading label"
		vars := mux.Vars(r)
		labelIDStr := vars["id"]

		id, err := entity.StringToID(labelIDStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid label ID"))
			return
		}

		tenant := r.Header.Get(httpHeaderTenantID)
		tenantID, err := entity.StringToID(tenant)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing tenant ID"))
			return
		}

		data, err := service.Get(tenantID, id)
		w.Header().Set(httpHeaderTenantID, tenant)
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

		toJ := &presenter.Label{
			ID:    data.ID,
			Name:  data.Name,
			Color: data.Color,
		}

		w.Header().Set(httpHeaderTenantID, data.TenantID.String())
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to encode label"))
		}
	})
}

func deleteLabel(service label.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error removing label"
		vars := mux.Vars(r)
		labelIDStr := vars["id"]

		id, err := entity.StringToID(labelIDStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid label ID"))
			return
		}

		tenant := r.Header.Get(httpHeaderTenantID)
		tenantID, err := entity.StringToID(tenant)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing tenant ID"))
			return
		}

		err = service.Delete(tenantID, id)
		w.Header().Set(httpHeaderTenantID, tenant)
		switch err {
		case nil:
			w.WriteHeader(http.StatusOK)
			return
		case entity.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Label doesn't exist"))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

func updateLabel(service label.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error updating label"

		vars := mux.Vars(r)
		labelIDStr := vars["id"]

		id, err := entity.StringToID(labelIDStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid label ID"))
			return
		}

		var input entity.Label
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
		err = service.Update(&input)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage + ":" + err.Error()))
			return
		}

		toJ := &presenter.Label{
			ID:       input.ID,
			TenantID: tenantID,
			Name:     input.Name,
			Color:    input.Color,
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

// MakeLabelHandlers make url handlers
func MakeLabelHandlers(r *mux.Router, n negroni.Negroni, service label.UseCase) {
	r.Handle("/v1/labels", n.With(
		negroni.Wrap(listLabels(service)),
	)).Methods("GET", "OPTIONS").Name("listLabels")

	r.Handle("/v1/labels", n.With(
		negroni.Wrap(createLabel(service)),
	)).Methods("POST", "OPTIONS").Name("createLabel")

	r.Handle("/v1/labels/{id}", n.With(
		negroni.Wrap(getLabel(service)),
	)).Methods("GET", "OPTIONS").Name("getLabel")

	r.Handle("/v1/labels/{id}", n.With(
		negroni.Wrap(deleteLabel(service)),
	)).Methods("DELETE", "OPTIONS").Name("deleteLabel")

	r.Handle("/v1/labels/{id}", n.With(
		negroni.Wrap(updateLabel(service)),
	)).Methods("PUT", "OPTIONS").Name("updateLabel")
}
