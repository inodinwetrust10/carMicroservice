// Package main
package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Starting API Gateway")

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.HandleFunc("POST /trip/preview", handleTripPreview)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("error starting the server")
	}
}
