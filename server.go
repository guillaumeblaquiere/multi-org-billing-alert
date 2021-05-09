package main

import (
	"fmt"
	"gblaquiere.dev/multi-org-billing-alert/handler"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/pubsub", handler.HandlePubsubMessage)
	http.HandleFunc("/http", handler.HandleHttpRequestMessage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
