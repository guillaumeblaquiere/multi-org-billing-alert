package main

import (
	"fmt"
	"gblaquiere.dev/multi-org-billing-alert/handler"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/pubsub", handler.HandlePubsubMessage).Methods(http.MethodPost)
	r.HandleFunc("/http", handler.UpsertBudgetAlert).Methods(http.MethodPost)
	r.HandleFunc("/http/projectid/{projectid}", handler.GetBudgetAlert).Methods(http.MethodGet)
	r.HandleFunc("/http/projectid/{projectid}", handler.DeleteBudgetAlert).Methods(http.MethodDelete)
	http.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
