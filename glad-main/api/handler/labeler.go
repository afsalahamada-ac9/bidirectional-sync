/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package handler

import (
	"log"
	"net/http"
	"sudhagar/glad/entity"
	"sudhagar/glad/pkg/common"
	"sudhagar/glad/usecase/labeler"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func setLabel(service labeler.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error setting label"
		var err error
		tenant := r.Header.Get(common.HttpHeaderTenantID)

		tenantID, err := entity.StringToID(tenant)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unable to parse tenant id"))
			return
		}

		vars := mux.Vars(r)
		labelID, err := entity.StringToID(vars["label_id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid label ID"))
			return
		}

		contactID, err := entity.StringToID(vars["contact_id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid contact ID"))
			return
		}

		err = service.SetLabel(tenantID, contactID, labelID)
		w.Header().Set(common.HttpHeaderTenantID, tenant)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func removeLabel(service labeler.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error removing label"
		var err error
		tenant := r.Header.Get(common.HttpHeaderTenantID)

		tenantID, err := entity.StringToID(tenant)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unable to parse tenant id"))
			return
		}

		vars := mux.Vars(r)
		labelID, err := entity.StringToID(vars["label_id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid label ID"))
			return
		}

		contactID, err := entity.StringToID(vars["contact_id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid contact ID"))
			return
		}

		err = service.RemoveLabel(tenantID, contactID, labelID)
		w.Header().Set(common.HttpHeaderTenantID, tenant)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// MakeLabelerHandlers make url handlers
func MakeLabelerHandlers(r *mux.Router, n negroni.Negroni, service labeler.UseCase) {

	r.Handle("/v1/contacts/label/{contact_id}/{label_id}", n.With(
		negroni.Wrap(setLabel(service)),
	)).Methods("PUT", "OPTIONS").Name("setLabel")

	r.Handle("/v1/contacts/label/{contact_id}/{label_id}", n.With(
		negroni.Wrap(removeLabel(service)),
	)).Methods("DELETE", "OPTIONS").Name("removeLabel")
}
