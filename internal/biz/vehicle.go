package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// Vehicle is a Vehicle model.
type Vehicle struct {
	// ID is the unique identifier of the vehicle.
	ID     int
	// VIN is the Vehicle Identification Number.
	VIN    string
	// UserID is the ID of the user who owns the vehicle.
	UserID int
}

// VehicleRepo defines the data access layer for Vehicle.
type VehicleRepo interface {
	// CreateVehicle creates a new vehicle.
	CreateVehicle(ctx context.Context, veh *Vehicle) error
	// FindOne finds a single vehicle by its ID.
	FindOne(ctx context.Context, id int) (*Vehicle, error)
	// FindByUserID finds all vehicles for a given user ID.
	FindByUserID(ctx context.Context, userID int) ([]*Vehicle, error)
}

// VehicleUsecase is a Vehicle usecase.
type VehicleUsecase struct {
	vehicleRepo VehicleRepo
	log         *log.Helper
}

// NewVehicleUsecase creates a Vehicle usecase.
func NewVehicleUsecase(vehicleRepo VehicleRepo, logger log.Logger) *VehicleUsecase {
	return &VehicleUsecase{vehicleRepo: vehicleRepo, log: log.NewHelper(logger)}
}
