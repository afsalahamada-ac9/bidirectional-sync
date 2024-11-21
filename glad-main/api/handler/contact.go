/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sudhagar/glad/pkg/common"
	"sudhagar/glad/usecase/contact"

	"sudhagar/glad/api/presenter"

	"sudhagar/glad/entity"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

const (
	apiContactListDefaultLimit int = 100
	apiContactListMaxLimit     int = 500
)

func listContacts(service contact.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading contacts"
		var data []*entity.Contact
		var err error
		tenant := r.Header.Get(common.HttpHeaderTenantID)

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
		limit := apiContactListDefaultLimit

		if qpIndex != "" {
			if index, err = strconv.Atoi(qpIndex); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Unable to parse index parameter"))
				return
			}
		}

		if qpLimit != "" {
			limit, err = strconv.Atoi(qpLimit)
			if err != nil || limit > apiContactListMaxLimit {
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
		w.Header().Set(common.HttpHeaderTenantID, tenant)

		var toJ []*presenter.Contact
		for _, d := range data {
			toJ = append(toJ, &presenter.Contact{
				ID:        d.ID,
				TenantID:  d.TenantID,
				AccountID: d.AccountID,
				Handle:    d.Handle,
				Name:      d.Name,
			})
		}
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.Header().Set(common.HttpHeaderTenantID, tenant)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to encode contact"))
		}
	})
}

// MakeContactHandlers make url handlers
func MakeContactHandlers(r *mux.Router, n negroni.Negroni, service contact.UseCase) {
	r.Handle("/v1/contacts", n.With(
		negroni.Wrap(listContacts(service)),
	)).Methods("GET", "OPTIONS").Name("listContacts")
}
