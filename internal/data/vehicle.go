package data

import (
	"context"
	"teslatrack/internal/biz"
)

var _ biz.VehicleRepo = (*vehicleRepo)(nil)

// vehicleRepo is the data layer implementation of VehicleRepo.
type vehicleRepo struct {
	data *Data
}

// NewVehicleRepo creates a new vehicleRepo.
func NewVehicleRepo(data *Data) biz.VehicleRepo {
	return &vehicleRepo{data}
}

// CreateVehicle implements biz.VehicleRepo.
func (v *vehicleRepo) CreateVehicle(ctx context.Context, veh *biz.Vehicle) error {
	panic("unimplemented")
}

// FindByUserID implements biz.VehicleRepo.
func (v *vehicleRepo) FindByUserID(ctx context.Context, userID int) ([]*biz.Vehicle, error) {
	panic("unimplemented")
}

// FindOne implements biz.VehicleRepo.
func (v *vehicleRepo) FindOne(ctx context.Context, id int) (*biz.Vehicle, error) {
	panic("unimplemented")
}
