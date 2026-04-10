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

	"ride-sharing/shared/env"
)

var httpAddr = env.GetString("HTTP_ADDR", ":8081")

func main() {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	log.Printf("Starting API Gateway on port %s:", httpAddr)
	mux.HandleFunc("POST /trip/preview", enableCORS(handleTripPreview))
	mux.HandleFunc("/ws/drivers", handleDriversWebSocket)
	mux.HandleFunc("/ws/riders", handleRidersWebSocket)

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
