package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type Vehicle struct {
	ID     int
	VIN    string
	UserID int
}

type VehicleRepo interface {
	CreateVehicle(ctx context.Context, veh *Vehicle) error
	FindOne(ctx context.Context, id int) (*Vehicle, error)
	FindByUserID(ctx context.Context, userID int) ([]*Vehicle, error)
}

type VehicleUsecase struct {
	vehicleRepo VehicleRepo
	log         *log.Helper
}

func NewVehicleUsecase(vehicleRepo VehicleRepo, logger log.Logger) *VehicleUsecase {
	return &VehicleUsecase{vehicleRepo: vehicleRepo, log: log.NewHelper(logger)}
}
