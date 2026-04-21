// Package service
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"ride-sharing/services/trip-service/internal/domain"
	tripType "ride-sharing/services/trip-service/pkg/types"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	repo domain.TripRepository
}

func NewService(repo domain.TripRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	t := &domain.TripModel{
		ID:       primitive.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: fare,
	}
	return s.repo.CreateTrip(ctx, t)
}

func (s *service) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*tripType.OsrmApiResponse, error) {
	url := fmt.Sprintf("http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson", pickup.Longitude, pickup.Latitude, destination.Longitude, destination.Latitude)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to route from the OSRM API %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response %v", err)
	}
	var routeResp tripType.OsrmApiResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, fmt.Errorf("failed to parse the response body %v", err)
	}
	return &routeResp, nil
}

func (s *service) EstimatePackagesPriceWithRoute(route *tripType.OsrmApiResponse) []*domain.RideFareModel {
	baseFares := getBaseFares()
	estimatedFares := make([]*domain.RideFareModel, len(baseFares))
	for i, f := range baseFares {
		estimatedFares[i] = estimateFareRoute(f, route)
	}
	return estimatedFares
}

func (s *service) GenerateTripFare(ctx context.Context, rideFares []*domain.RideFareModel, userID string) ([]*domain.RideFareModel, error) {
	fares := make([]*domain.RideFareModel, len(rideFares))
	for i, f := range rideFares {
		id := primitive.NewObjectID()
		fare := &domain.RideFareModel{
			ID:          id,
			UserID:      userID,
			TotalPrice:  f.TotalPrice,
			PackageSlug: f.PackageSlug,
		}
		if err := s.repo.SaveRideFare(ctx, fare); err != nil {
			return nil, fmt.Errorf("failed to save trip fare %v", err)
		}
		fares[i] = fare
	}
	return fares, nil
}

func estimateFareRoute(fare *domain.RideFareModel, route *tripType.OsrmApiResponse) *domain.RideFareModel {
	pricingCfg := tripType.DefaultPricingConfig()
	carPackagePrice := fare.TotalPrice

	distanceKm := route.Routes[0].Distance
	durationMin := route.Routes[0].Duration

	distanceFare := distanceKm * pricingCfg.PricePerDis
	timeFare := durationMin * pricingCfg.PricePerMin

	totalPrice := carPackagePrice + distanceFare + timeFare

	return &domain.RideFareModel{
		TotalPrice:  totalPrice,
		PackageSlug: fare.PackageSlug,
	}
}

func getBaseFares() []*domain.RideFareModel {
	return []*domain.RideFareModel{
		{
			PackageSlug: "suv",
			TotalPrice:  200,
		},
		{
			PackageSlug: "sedan",
			TotalPrice:  350,
		},
		{
			PackageSlug: "van",
			TotalPrice:  400,
		},
		{
			PackageSlug: "luxury",
			TotalPrice:  1000,
		},
	}
}
