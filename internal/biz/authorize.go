package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// Authorize is the model for authorization.
type Authorize struct {
	ID           int64
	ClientID     string
	ClientSecret string
	GrantType    string
	RedirectURI  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// AuthorizeRepo defines the storage interface for Authorize.
type AuthorizeRepo interface {
	Create(context.Context, *Authorize) (*Authorize, error)
	Update(context.Context, *Authorize) (*Authorize, error)
	FindByID(context.Context, int64) (*Authorize, error)
}

// AuthorizeUsecase is the use case for authorization.
type AuthorizeUsecase struct {
	repo AuthorizeRepo
	log  *log.Helper
}

// NewAuthorizeUsecase creates a new AuthorizeUsecase.
func NewAuthorizeUsecase(repo AuthorizeRepo, logger log.Logger) *AuthorizeUsecase {
	return &AuthorizeUsecase{repo: repo, log: log.NewHelper(logger)}
}
