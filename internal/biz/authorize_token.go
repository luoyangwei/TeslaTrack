package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// AuthorizeToken is the business model for Tesla API authorization tokens.
// It holds the token information needed to make authenticated API calls.
type AuthorizeToken struct {
	ID           int64     // Unique identifier for the token record.
	TeslaCode    string    // The authorization code from Tesla's OAuth flow.
	ClientID     string    // The client ID used for the Tesla API.
	ClientSecret string    // The client secret used for the Tesla API.
	AccessToken  string    // The token used to access the Tesla API.
	RefreshToken string    // The token used to refresh the access token.
	Scope        string    // The scope of permissions granted.
	CreatedAt    time.Time // The timestamp when the token was created.
	UpdatedAt    time.Time // The timestamp when the token was last updated.
	Deleted      bool      // A flag for soft deletion.
}

// AuthorizeTokenRepo defines the persistence layer interface for AuthorizeToken data.
type AuthorizeTokenRepo interface {
	// Create saves a new AuthorizeToken record.
	Create(ctx context.Context, token *AuthorizeToken) (*AuthorizeToken, error)
	// Update modifies an existing AuthorizeToken record.
	Update(ctx context.Context, token *AuthorizeToken) error
	// FindByClientID retrieves an AuthorizeToken by the client ID.
	FindByClientID(ctx context.Context, clientID string) (*AuthorizeToken, error)
	// FindByAccessToken retrieves an AuthorizeToken by the access token.
	FindByAccessToken(ctx context.Context, accessToken string) (*AuthorizeToken, error)
	// Delete soft-deletes an AuthorizeToken record by its ID.
	Delete(ctx context.Context, id int64) error
}

// AuthorizeTokenUsecase provides the business logic for authorization token operations.
type AuthorizeTokenUsecase struct {
	repo AuthorizeTokenRepo
	log  *log.Helper
}

// NewAuthorizeTokenUsecase creates a new instance of AuthorizeTokenUsecase.
func NewAuthorizeTokenUsecase(repo AuthorizeTokenRepo, logger log.Logger) *AuthorizeTokenUsecase {
	return &AuthorizeTokenUsecase{repo: repo, log: log.NewHelper(logger)}
}

// Create is the use case for creating a new authorization token.
func (uc *AuthorizeTokenUsecase) Create(ctx context.Context, token *AuthorizeToken) (*AuthorizeToken, error) {
	return uc.repo.Create(ctx, token)
}

// Update is the use case for updating an existing authorization token.
func (uc *AuthorizeTokenUsecase) Update(ctx context.Context, token *AuthorizeToken) error {
	return uc.repo.Update(ctx, token)
}

// FindByClientID is the use case for finding a token by client ID.
func (uc *AuthorizeTokenUsecase) FindByClientID(ctx context.Context, clientID string) (*AuthorizeToken, error) {
	return uc.repo.FindByClientID(ctx, clientID)
}

// FindByAccessToken is the use case for finding a token by access token.
func (uc *AuthorizeTokenUsecase) FindByAccessToken(ctx context.Context, accessToken string) (*AuthorizeToken, error) {
	return uc.repo.FindByAccessToken(ctx, accessToken)
}

// Delete is the use case for deleting (soft delete) an authorization token.
func (uc *AuthorizeTokenUsecase) Delete(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}
