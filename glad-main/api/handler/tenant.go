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

	"sudhagar/glad/usecase/tenant"

	"sudhagar/glad/api/presenter"

	"sudhagar/glad/entity"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func listTenants(service tenant.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading tenants"
		var data []*entity.Tenant
		var err error

		data, err = service.ListTenants()
		w.Header().Set("Content-Type", "application/json")
		if err != nil && err != entity.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage + ":" + err.Error()))
			return
		}

		total := service.GetCount()
		w.Header().Set(httpHeaderTotalCount, strconv.Itoa(total))

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}

		var toJ []*presenter.Tenant
		for _, d := range data {
			toJ = append(toJ, &presenter.Tenant{
				ID:        d.ID,
				Name:      d.Name,
				Country:   d.Country,
				AuthToken: d.AuthToken,
			})
		}
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to encode tenant"))
		}
	})
}

func createTenant(service tenant.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error adding tenant"
		var input struct {
			Name    string `json:"name"`
			Country string `json:"country"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unable to decode the data. " + err.Error()))
			return
		}

		id, err := service.CreateTenant(input.Name, input.Country)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage + ":" + err.Error()))
			return
		}
		toJ := &presenter.Tenant{
			ID:      id,
			Name:    input.Name,
			Country: input.Country,
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

func getTenant(service tenant.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading tenant"
		vars := mux.Vars(r)
		id, err := entity.StringToID(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		data, err := service.GetTenant(id)
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

		toJ := &presenter.Tenant{
			ID:        data.ID,
			Name:      data.Name,
			Country:   data.Country,
			AuthToken: data.AuthToken,
		}

		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to encode tenant"))
		}
	})
}

func deleteTenant(service tenant.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error removing tenant"
		vars := mux.Vars(r)
		id, err := entity.StringToID(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		err = service.DeleteTenant(id)
		switch err {
		case nil:
			w.WriteHeader(http.StatusOK)
			return
		case entity.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Tenant doesn't exist"))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

func login(service tenant.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading tenant"
		var input struct {
			Name    string `json:"name"`
			Country string `json:"country"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unable to decode the data. " + err.Error()))
			return
		}

		data, err := service.Login(input.Name, input.Country)
		switch err {
		case nil:
			break
		// intentionally returning same response for auth failure and not found scenarios
		case entity.ErrAuthFailure, entity.ErrNotFound:
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid login credentials"))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage + ":" + err.Error()))
			return
		}

		toJ := &presenter.Tenant{
			ID:        data.ID,
			Name:      data.Name,
			AuthToken: data.AuthToken,
		}

		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to encode tenant"))
		}
	})
}

func updateTenant(service tenant.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error updating tenant"

		vars := mux.Vars(r)
		id, err := entity.StringToID(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		var input entity.Tenant
		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unable to decode the data. " + err.Error()))
			return
		}

		err = service.UpdateTenant(&input)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage + ":" + err.Error()))
			return
		}

		input.ID = id
		toJ := &presenter.Tenant{
			ID:      input.ID,
			Name:    input.Name,
			Country: input.Country,
			// AuthToken is not returned back
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

// MakeTenantHandlers make url handlers
func MakeTenantHandlers(r *mux.Router, n negroni.Negroni, service tenant.UseCase) {
	r.Handle("/v1/tenants", n.With(
		negroni.Wrap(listTenants(service)),
	)).Methods("GET", "OPTIONS").Name("listTenants")

	r.Handle("/v1/tenants", n.With(
		negroni.Wrap(createTenant(service)),
	)).Methods("POST", "OPTIONS").Name("createTenant")

	r.Handle("/v1/tenants/{id}", n.With(
		negroni.Wrap(getTenant(service)),
	)).Methods("GET", "OPTIONS").Name("getTenant")

	r.Handle("/v1/tenants/{id}", n.With(
		negroni.Wrap(deleteTenant(service)),
	)).Methods("DELETE", "OPTIONS").Name("deleteTenant")

	r.Handle("/v1/login", n.With(
		negroni.Wrap(login(service)),
	)).Methods("POST", "OPTIONS").Name("login")

	r.Handle("/v1/tenants/{id}", n.With(
		negroni.Wrap(updateTenant(service)),
	)).Methods("PUT", "OPTIONS").Name("updateTenant")
}
