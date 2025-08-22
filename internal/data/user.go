package data

import (
	"teslatrack/internal/biz"
)

var _ biz.UserRepo = (*userRepo)(nil)

// userRepo is the data layer implementation of UserRepo.
type userRepo struct {
	data *Data
}

// NewUserRepo creates a new userRepo.
func NewUserRepo(data *Data) biz.UserRepo {
	return &userRepo{data}
}
