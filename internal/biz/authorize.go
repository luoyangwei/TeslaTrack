package biz

import (
	"context"
	"teslatrack/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

// ALL_SCOPES defines the full list of permissions the application can request from Tesla's API.
// - openid: Allows login with Tesla credentials.
// - offline_access: Allows obtaining a refresh token for offline access.
// - user_data: Access to user profile information.
// - vehicle_device_data: Access to vehicle information and real-time data.
// - vehicle_location: Access to vehicle location information.
// - vehicle_cmds: Allows sending commands to the vehicle (e.g., unlock, start).
// - vehicle_charging_cmds: Allows managing vehicle charging.
// - energy_device_data: Access to energy product information (e.g., Powerwall).
// - energy_cmds: Allows sending commands to energy products.
const ALL_SCOPES = "openid offline_access user_data vehicle_device_data vehicle_location vehicle_cmds vehicle_charging_cmds energy_device_data energy_cmds"

// TESLA_EXCHANGE_CODE_URL
const TESLA_EXCHANGE_CODE_URL = "https://auth.tesla.cn/oauth2/v3/token"

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
	conf *conf.Server
	log  *log.Helper
}

// NewAuthorizeUsecase creates a new instance of AuthorizeUsecase.
// It requires an AuthorizeRepo for data access and a logger for logging.
func NewAuthorizeUsecase(repo AuthorizeRepo, config *conf.Server, logger log.Logger) *AuthorizeUsecase {
	return &AuthorizeUsecase{repo: repo, conf: config, log: log.NewHelper(logger)}
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

// Callback handles the authorization code received from the OAuth provider.
// It is responsible for exchanging the code for an access token.
// NOTE: This method is a stub and its implementation is pending.
func (uc *AuthorizeUsecase) Callback(ctx context.Context, code string) error {
	return nil
}

// AuthorizeEncodeRedirect holds the parameters needed to construct the redirect URL
// for the Tesla OAuth 2.0 authorization flow.
type AuthorizeEncodeRedirect struct {
	ClientID               string `json:"clientId"`
	RedirectURI            string `json:"redirectUri"`
	Scope                  string `json:"scope"`
	State                  string `json:"state"`
	Nonce                  string `json:"nonce"`
	PromptMissingScopes    bool   `json:"promptMissingScopes"`
	RequireRequestedScopes bool   `json:"requireRequestedScopes"`
}

// Redirect prepares the necessary parameters for the authorization redirect.
// It fetches client details, generates a secure state and nonce, and returns them.
func (uc *AuthorizeUsecase) Redirect(ctx context.Context, clientID string) (*AuthorizeEncodeRedirect, error) {
	// Fetch the client's authorization configuration.
	authorize, err := uc.repo.FindByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Generate a unique and non-guessable value for state and nonce to prevent CSRF and replay attacks.
	state, _ := uuid.NewV7()
	nonce, _ := uuid.NewV7()

	// Construct the redirect parameters.
	return &AuthorizeEncodeRedirect{
		ClientID:               authorize.ClientID,
		RedirectURI:            authorize.RedirectURI,
		Scope:                  ALL_SCOPES,
		State:                  state.String(),
		Nonce:                  nonce.String(),
		PromptMissingScopes:    false, // These could be configurable in the future.
		RequireRequestedScopes: false,
	}, nil
}
