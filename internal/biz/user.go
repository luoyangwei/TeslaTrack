package biz

import (
	"github.com/go-kratos/kratos/v2/log"
)

// User is a User model.
type User struct {
	// ID is the unique identifier of the user.
	ID       int
	// Account is the user's account name.
	Account  string
	// Password is the user's password (hashed).
	Password string
	// Vehicles is the list of vehicles associated with the user.
	Vehicles []*Vehicle
}

// UserRepo defines the data access layer for User.
// It is currently a placeholder and should be expanded with methods.
type UserRepo interface{}

// UserUsecase is a User usecase.
type UserUsecase struct {
	userRepo UserRepo
	log      *log.Helper
}

// NewUserUsecase creates a User usecase.
func NewUserUsecase(userRepo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{userRepo: userRepo, log: log.NewHelper(logger)}
}
