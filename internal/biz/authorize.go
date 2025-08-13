package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// Authorize is the data model for OAuth 2.0 client authorization.
// It holds the necessary information for a client to obtain an access token.
type Authorize struct {
	ID           int64     // Unique identifier for the authorization record.
	ClientID     string    // The client identifier, issued to the client during the registration process.
	ClientSecret string    // The client secret, a confidential value used to authenticate the client.
	GrantType    string    // The grant type the client is allowed to use (e.g., "authorization_code", "client_credentials").
	RedirectURI  string    // The URI to which the authorization server will redirect the user-agent after granting authorization.
	CreatedAt    time.Time // The timestamp when the authorization record was created.
	UpdatedAt    time.Time // The timestamp when the authorization record was last updated.
}

// AuthorizeRepo defines the persistence layer interface for Authorize data.
// This interface abstracts the underlying data storage (e.g., database, cache).
type AuthorizeRepo interface {
	// Create saves a new Authorize record to the storage.
	Create(ctx context.Context, authorize *Authorize) error
	// Update modifies an existing Authorize record in the storage.
	Update(ctx context.Context, authorize *Authorize) error
	// FindByClientID retrieves an Authorize record from the storage by its client ID.
	FindByClientID(ctx context.Context, clientID string) (*Authorize, error)
}

// AuthorizeUsecase provides the business logic for authorization operations.
// It orchestrates the interaction between the transport layer (e.g., HTTP server) and the data layer (repository).
type AuthorizeUsecase struct {
	repo AuthorizeRepo
	log  *log.Helper
}

// NewAuthorizeUsecase creates a new instance of AuthorizeUsecase.
// It requires an AuthorizeRepo for data access and a logger for logging.
func NewAuthorizeUsecase(repo AuthorizeRepo, logger log.Logger) *AuthorizeUsecase {
	return &AuthorizeUsecase{repo: repo, log: log.NewHelper(logger)}
}

// Create is the use case for creating a new authorization.
// It delegates the creation operation to the underlying repository.
func (uc *AuthorizeUsecase) Create(ctx context.Context, authorize *Authorize) error {
	return uc.repo.Create(ctx, authorize)
}

// Update is the use case for updating an existing authorization.
// It delegates the update operation to the underlying repository.
func (uc *AuthorizeUsecase) Update(ctx context.Context, authorize *Authorize) error {
	return uc.repo.Update(ctx, authorize)
}

// FindByClientID is the use case for finding an authorization by its client ID.
// It delegates the find operation to the underlying repository.
func (uc *AuthorizeUsecase) FindByClientID(ctx context.Context, clientID string) (*Authorize, error) {
	return uc.repo.FindByClientID(ctx, clientID)
}
