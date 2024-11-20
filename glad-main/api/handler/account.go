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

	"sudhagar/glad/usecase/account"

	"sudhagar/glad/api/presenter"

	"sudhagar/glad/entity"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func listAccounts(service account.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading accounts"
		var data []*entity.Account
		var err error
		tenant := r.Header.Get(httpHeaderTenantID)

		tenantID, err := entity.StringToID(tenant)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unable to parse tenant id"))
			return
		}

		data, err = service.ListAccounts(tenantID)
		w.Header().Set("Content-Type", "application/json")
		if err != nil && err != entity.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage + ":" + err.Error()))
			return
		}

		total := service.GetCount(tenantID)
		w.Header().Set(httpHeaderTotalCount, strconv.Itoa(total))
		w.Header().Set(httpHeaderTenantID, tenant)

		var toJ []*presenter.Account
		for _, d := range data {
			toJ = append(toJ, &presenter.Account{
				ID:        d.ID,
				ExtID:     d.ExtID,
				Username:  d.Username,
				FirstName: d.FirstName,
				LastName:  d.LastName,
				Phone:     d.Phone,
				Email:     d.Email,
				Type:      d.Type,
			})
		}
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to encode account"))
		}
	})
}

// func createAccount(service account.UseCase) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		errorMessage := "Error adding account"
// 		var input struct {
// 			Name    string             `json:"name"`
// 			Type    entity.AccountType `json:"type"`
// 			Content string             `json:"content"`
// 		}

// 		tenant := r.Header.Get(httpHeaderTenantID)
// 		tenantID, err := entity.StringToID(tenant)
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			w.Write([]byte("Missing tenant ID"))
// 			return
// 		}

// 		err = json.NewDecoder(r.Body).Decode(&input)
// 		if err != nil {
// 			log.Println(err.Error())
// 			w.WriteHeader(http.StatusBadRequest)
// 			w.Write([]byte("Unable to decode the data. " + err.Error()))
// 			return
// 		}

// 		id, err := service.CreateAccount(
// 			tenantID,
// 			input.Name,
// 			input.Type,
// 			input.Content)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			w.Write([]byte(errorMessage + ":" + err.Error()))
// 			return
// 		}
// 		toJ := &presenter.Account{
// 			ID:       id,
// 			Name:     input.Name,
// 			Type:     entity.AccountText,
// 			Content:  input.Content,
// 			TenantID: tenantID,
// 		}

// 		w.Header().Set(httpHeaderTenantID, tenant)
// 		w.WriteHeader(http.StatusCreated)
// 		if err := json.NewEncoder(w).Encode(toJ); err != nil {
// 			log.Println(err.Error())
// 			w.WriteHeader(http.StatusInternalServerError)
// 			w.Write([]byte(errorMessage))
// 			return
// 		}
// 	})
// }

func getAccount(service account.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading account"
		vars := mux.Vars(r)
		username := vars["username"]
		data, err := service.GetAccount(username)
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

		toJ := &presenter.Account{
			ID:        data.ID,
			ExtID:     data.ExtID,
			Username:  data.Username,
			FirstName: data.FirstName,
			LastName:  data.LastName,
			Phone:     data.Phone,
			Email:     data.Email,
			Type:      data.Type,
		}

		w.Header().Set(httpHeaderTenantID, r.Header.Get(httpHeaderTenantID))
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to encode account"))
		}
	})
}

func deleteAccount(service account.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error removing account"
		vars := mux.Vars(r)
		username := vars["username"]
		err := service.DeleteAccount(username)
		switch err {
		case nil:
			w.WriteHeader(http.StatusOK)
			return
		case entity.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Account doesn't exist"))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

func updateAccount(service account.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error updating account"

		vars := mux.Vars(r)
		username := vars["username"]

		var input entity.Account
		// tenant := r.Header.Get(httpHeaderTenantID)
		// tenantID, err := entity.StringToID(tenant)
		// if err != nil {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	w.Write([]byte("Missing tenant ID"))
		// 	return
		// }

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unable to decode the data. " + err.Error()))
			return
		}

		input.Username = username
		err = service.UpdateAccount(&input)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage + ":" + err.Error()))
			return
		}

		toJ := &presenter.Account{
			ID:        input.ID,
			ExtID:     input.ExtID,
			Username:  input.Username,
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Phone:     input.Phone,
			Email:     input.Email,
			Type:      input.Type,
		}

		w.Header().Set(httpHeaderTenantID, "")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

// MakeAccountHandlers make url handlers
func MakeAccountHandlers(r *mux.Router, n negroni.Negroni, service account.UseCase) {
	r.Handle("/v1/accounts", n.With(
		negroni.Wrap(listAccounts(service)),
	)).Methods("GET", "OPTIONS").Name("listAccounts")

	// r.Handle("/v1/accounts", n.With(
	// 	negroni.Wrap(createAccount(service)),
	// )).Methods("POST", "OPTIONS").Name("createAccount")

	r.Handle("/v1/accounts/{username}", n.With(
		negroni.Wrap(getAccount(service)),
	)).Methods("GET", "OPTIONS").Name("getAccount")

	r.Handle("/v1/accounts/{username}", n.With(
		negroni.Wrap(deleteAccount(service)),
	)).Methods("DELETE", "OPTIONS").Name("deleteAccount")

	r.Handle("/v1/accounts/{username}", n.With(
		negroni.Wrap(updateAccount(service)),
	)).Methods("PUT", "OPTIONS").Name("updateAccount")
}
