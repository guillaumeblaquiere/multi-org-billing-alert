package handler

import (
	"encoding/json"
	"fmt"
	"gblaquiere.dev/multi-org-billing-alert/internal"
	"gblaquiere.dev/multi-org-billing-alert/internal/httperrors"
	"gblaquiere.dev/multi-org-billing-alert/model"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

/*
	Received the JSON message of Alert Creation in HTTP Request
*/

var projectIdParam = "projectid"

func DeleteBudgetAlert(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	projectId := vars[projectIdParam]

	err := internal.DeleteBillingAlert(r.Context(), projectId)
	if err != nil {
		log.Printf("internal.DeleteBillingAlert: %v\n", err)
		http.Error(w, err.Error(), httperrors.GetHttpCode(err))
		return
	}
}

func GetBudgetAlert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectId := vars[projectIdParam]

	billingAlert, err := internal.GetBillingAlert(r.Context(), projectId)
	if err != nil {
		log.Printf("internal.GetBillingAlert: %v\n", err)
		http.Error(w, err.Error(), httperrors.GetHttpCode(err))
		return
	}

	json, err := json.Marshal(billingAlert)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func UpsertBudgetAlert(w http.ResponseWriter, r *http.Request) {

	var billing model.BillingAlert
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %v\n", err)
		http.Error(w, fmt.Sprintf("Bad Request %q", err), http.StatusBadRequest)
		return
	}
	log.Printf("PubSub message content:\n %s\n", string(body))

	if err := json.Unmarshal(body, &billing); err != nil {
		log.Printf("json.Unmarshal: %v\n", err)
		http.Error(w, fmt.Sprintf("Bad Request %q", err), http.StatusBadRequest)
		return
	}
	err = internal.CreateBillingAlert(r.Context(), &billing)
	if err != nil {
		log.Printf("internal.CreateBillingAlert: %v\n", err)
		http.Error(w, err.Error(), httperrors.GetHttpCode(err))
		return
	}
	w.WriteHeader(http.StatusCreated)
}
