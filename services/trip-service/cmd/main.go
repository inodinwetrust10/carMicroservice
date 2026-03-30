package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	h "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
)

func main() {
	inmemRepo := repository.NewInmemRepository()
	svc := service.NewService(inmemRepo)
	mux := http.NewServeMux()
	httphandler := h.HttpHandler{Service: svc}
	mux.HandleFunc("POST /preview", httphandler.HandleTripPreview)

	server := &http.Server{
		Handler: mux,
		Addr:    ":8083",
	}
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
