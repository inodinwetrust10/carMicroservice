// Package domain
package domain

import (
	"context"

	tripType "ride-sharing/services/trip-service/pkg/types"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripModel struct {
	ID       primitive.ObjectID
	UserID   string
	Status   string
	RideFare *RideFareModel
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *TripModel) (*TripModel, error)
	SaveRideFare(ctx context.Context, f *RideFareModel) error
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*tripType.OsrmApiResponse, error)
	EstimatePackagesPriceWithRoute(route *tripType.OsrmApiResponse) []*RideFareModel
	GenerateTripFare(ctx context.Context, fares []*RideFareModel, userID string) ([]*RideFareModel, error)
}
