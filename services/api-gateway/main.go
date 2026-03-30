// Package main
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Starting API Gateway")

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.HandleFunc("POST /trip/preview", handleTripPreview)

	done := make(chan os.Signal, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("error starting the server")
		}
	}()

	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	<-done
	log.Println("shutting down the server....")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		if err := server.Close(); err != nil {
			log.Printf("forced shutdown failed: %v", err)
		}
	}
}
