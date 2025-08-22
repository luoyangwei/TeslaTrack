package biz

import (
	"github.com/go-kratos/kratos/v2/log"
)

type User struct {
	ID       int
	Account  string
	Password string
	Vehicles []*Vehicle
}

type UserRepo interface{}

type UserUsecase struct {
	userRepo UserRepo
	log      *log.Helper
}

func NewUserUsecase(userRepo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{userRepo: userRepo, log: log.NewHelper(logger)}
}
