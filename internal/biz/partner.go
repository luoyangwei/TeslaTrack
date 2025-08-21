package biz

import (
	"context"
	"os"
	"teslatrack/internal/conf"
	"teslatrack/pkg/tesla"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// Partner is a Partner model.
type Partner struct {
	// ID is the id of the partner.
	ID int
	// ClientID is the client id of the partner.
	ClientID string
	// AccessToken is the access token of the partner.
	AccessToken string
	// ExpiresIn is the expires in of the partner.
	ExpiresIn int32
	// TokenType is the token type of the partner.
	TokenType string
	// CreatedAt is the created at of the partner.
	CreatedAt time.Time
	// UpdatedAt is the updated at of the partner.
	UpdatedAt time.Time
}

// PartnerRepo is a Partner repo.
type PartnerRepo interface {
	// Get gets a Partner by clientID.
	Get(ctx context.Context, clientID string) (*Partner, error)
	// MustGet gets a Partner by clientID, returns nil if not found.
	MustGet(ctx context.Context, clientID string) (*Partner, error)
	// Create creates a Partner.
	Create(ctx context.Context, partner *Partner) error
	// Update updates a Partner.
	Update(ctx context.Context, id int, partner *Partner) error
}

// PartnerUsecase is a Partner usecase.
type PartnerUsecase struct {
	repo PartnerRepo
	conf *conf.Server
	log  *log.Helper
}

// NewPartnerUsecase creates a Partner usecase.
func NewPartnerUsecase(repo PartnerRepo, conf *conf.Server, logger log.Logger) *PartnerUsecase {
	return &PartnerUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

// Initialize is server starting initialize.
// It ensures that the partner information from Tesla is stored in the database.
func (uc *PartnerUsecase) Initialize() error {
	ctx := context.Background()
	// Get Tesla client credentials from environment variables.
	clientID, clientSecret := os.Getenv("TESLA_CLIENT_ID"), os.Getenv("TESLA_CLIENT_SECRET")
	// Try to get the partner from the database.
	partner, err := uc.repo.MustGet(ctx, clientID)
	if err != nil {
		return err
	}

	// If partner is not found in the database, fetch from Tesla API and create it.
	if partner == nil {
		uc.log.Info("Partner not found in database, fetching from Tesla API.")
		teslaPartner, err := tesla.GetPartner(clientID, clientSecret)
		if err != nil {
			return err
		}

		uc.log.Infow("msg", "Create the partner in the database.", "access_token", teslaPartner.AccessToken, "expires_in", teslaPartner.ExpiresIn)

		// Create the partner in the database.
		partner = &Partner{
			ClientID:    clientID,
			AccessToken: teslaPartner.AccessToken,
			ExpiresIn:   int32(teslaPartner.ExpiresIn),
			TokenType:   teslaPartner.TokenType,
			CreatedAt:   time.Now(),
		}
		err = uc.repo.Create(ctx, partner)
		if err != nil {
			return err
		}
	}

	// Check if the partner token has expired.
	if partner.CreatedAt.Add(time.Duration(partner.ExpiresIn) * time.Second).Before(time.Now()) {
		uc.log.Info("Partner token expired, refreshing...")
		// Fetch a new token from Tesla API.
		teslaPartner, err := tesla.GetPartner(clientID, clientSecret)
		if err != nil {
			return err
		}

		// Update partner details with the new token.
		partner.AccessToken = teslaPartner.AccessToken
		partner.ExpiresIn = int32(teslaPartner.ExpiresIn)
		partner.TokenType = teslaPartner.TokenType

		// Persist the updated partner info.
		err = uc.repo.Update(ctx, partner.ID, partner)
		if err != nil {
			return err
		}
		uc.log.Infow("msg", "Partner token refreshed successfully.", "access_token", partner.AccessToken, "expires_in", partner.ExpiresIn)
	}

	uc.log.Infow("msg", "Partner initialized successfully.", "access_token", partner.AccessToken, "expires_in", partner.ExpiresIn)
	return nil
}
