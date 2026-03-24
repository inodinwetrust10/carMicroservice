package main

import (
	"context"
	"log"

	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
)

func main() {
	inmemRepo := repository.NewInmemRepository()
	svc := service.NewService(inmemRepo)
	fare := &domain.RideFareModel{
		UserID: "42",
	}
	t, err := svc.CreateTrip(context.Background(), fare)
	if err != nil {
		log.Println(err)
	}
	log.Println(t)
}
