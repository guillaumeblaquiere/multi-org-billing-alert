package handler

import (
	"encoding/json"
	"errors"
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
var alertNameParam = "alertname"

func DeleteBudgetAlert(w http.ResponseWriter, r *http.Request) {

	name, err := getAlertName(r)

	if err != nil {
		http.Error(w, err.Error(), httperrors.GetHttpCode(err))
		return
	}

	billingAlert, err := internal.DeleteBillingAlert(r.Context(), name)
	if err != nil {
		log.Printf("internal.DeleteBillingAlert: %v\n", err)
		http.Error(w, err.Error(), httperrors.GetHttpCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	formatResponse(w, billingAlert)
}

func GetBudgetAlert(w http.ResponseWriter, r *http.Request) {

	name, err := getAlertName(r)

	if err != nil {
		http.Error(w, err.Error(), httperrors.GetHttpCode(err))
		return
	}

	billingAlert, err := internal.GetBillingAlert(r.Context(), name)
	if err != nil {
		log.Printf("internal.GetBillingAlert: %v\n", err)
		http.Error(w, err.Error(), httperrors.GetHttpCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	formatResponse(w, billingAlert)
}

func getAlertName(r *http.Request) (name string, err error) {
	vars := mux.Vars(r)
	projectId := vars[projectIdParam]
	alertName := vars[alertNameParam]

	if projectId != "" && alertName != "" && projectId != alertName {
		return "", httperrors.New(errors.New("projectID and AlertName provided and different"), http.StatusBadRequest)
	}

	name = projectId
	if alertName != "" {
		name = alertName
	}
	return
}

func UpsertBudgetAlert(w http.ResponseWriter, r *http.Request) {

	var billing model.BillingAlert
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %v\n", err)
		http.Error(w, fmt.Sprintf("Bad Request %q", err), http.StatusBadRequest)
		return
	}
	log.Printf("Message content:\n %s\n", string(body))

	if err := json.Unmarshal(body, &billing); err != nil {
		log.Printf("json.Unmarshal: %v\n", err)
		http.Error(w, fmt.Sprintf("Bad Request %q", err), http.StatusBadRequest)
		return
	}

	billingAlert, err := internal.CreateBillingAlert(r.Context(), &billing)
	if err != nil {
		log.Printf("internal.CreateBillingAlert: %v\n", err)
		http.Error(w, err.Error(), httperrors.GetHttpCode(err))
		return
	}
	w.WriteHeader(http.StatusCreated)
	formatResponse(w, billingAlert)
}

func formatResponse(w http.ResponseWriter, billingAlert *model.BillingAlert) {
	billingAlertJson, _ := json.Marshal(billingAlert)
	fmt.Fprint(w, string(billingAlertJson))
	w.Header().Add("Content-type", "application/json")
}
