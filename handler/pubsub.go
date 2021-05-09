package handler

import (
	"encoding/json"
	"gblaquiere.dev/multi-org-billing-alert/internal"
	"gblaquiere.dev/multi-org-billing-alert/model"
	"io/ioutil"
	"log"
	"net/http"
)

/*
	Received the JSON message of Alert Creation in PubSub message
*/
func HandlePubsubMessage(w http.ResponseWriter, r *http.Request) {
	var m model.PubSubMessage
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %v\n", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &m); err != nil {
		log.Printf("json.Unmarshal: %v\n", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var billing model.BillingAlert
	log.Printf("PubSub message content:\n %s\n", string(m.Message.Data))

	if err := json.Unmarshal(m.Message.Data, &billing); err != nil {
		log.Printf("json.Unmarshal: %v\n", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = internal.CreateBillingAlert(r.Context(), &billing)
	if err != nil {
		log.Printf("internal.CreateBillingAlert: %v\n", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
}
