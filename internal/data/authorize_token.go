package data

import (
	"context"
	"teslatrack/internal/biz"
	"teslatrack/internal/data/ent"
	"teslatrack/internal/data/ent/authorizetoken"
	"time"
)

// A compile-time check to ensure that authorizeTokenRepo implements the biz.AuthorizeTokenRepo interface.
var _ biz.AuthorizeTokenRepo = (*authorizeTokenRepo)(nil)

// authorizeTokenRepo is the data access layer implementation for authorization tokens.
// It uses the ent framework to interact with the database.
type authorizeTokenRepo struct {
	data *Data
}

// NewAuthorizeTokenRepo creates a new authorizeTokenRepo.
func NewAuthorizeTokenRepo(data *Data) biz.AuthorizeTokenRepo {
	return &authorizeTokenRepo{data: data}
}

// toBizToken converts an ent.AuthorizeToken model to a biz.AuthorizeToken model.
func toBizToken(model *ent.AuthorizeToken) *biz.AuthorizeToken {
	if model == nil {
		return nil
	}
	return &biz.AuthorizeToken{
		ID:           int64(model.ID),
		TeslaCode:    model.TeslaCode,
		ClientID:     model.ClientID,
		ClientSecret: model.ClientSecret,
		AccessToken:  model.AccessToken,
		RefreshToken: model.RefreshToken,
		Scope:        model.Scope,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
		Deleted:      model.Deleted,
	}
}

// Create saves a new authorization token record to the database.
func (r *authorizeTokenRepo) Create(ctx context.Context, token *biz.AuthorizeToken) (*biz.AuthorizeToken, error) {
	model, err := r.data.db.AuthorizeToken.Create().
		SetTeslaCode(token.TeslaCode).
		SetClientID(token.ClientID).
		SetClientSecret(token.ClientSecret).
		SetAccessToken(token.AccessToken).
		SetRefreshToken(token.RefreshToken).
		SetScope(token.Scope).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return toBizToken(model), nil
}

// Update modifies an existing authorization token record in the database.
func (r *authorizeTokenRepo) Update(ctx context.Context, token *biz.AuthorizeToken) error {
	_, err := r.data.db.AuthorizeToken.UpdateOneID(int(token.ID)).
		SetAccessToken(token.AccessToken).
		SetRefreshToken(token.RefreshToken).
		SetScope(token.Scope).
		SetUpdatedAt(time.Now()). // Explicitly update the timestamp
		Save(ctx)
	return err
}

// FindByClientID retrieves a token by its client ID.
func (r *authorizeTokenRepo) FindByClientID(ctx context.Context, clientID string) (*biz.AuthorizeToken, error) {
	model, err := r.data.db.AuthorizeToken.Query().
		Where(authorizetoken.ClientID(clientID), authorizetoken.Deleted(false)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return toBizToken(model), nil
}

// FindByAccessToken retrieves a token by its access token.
func (r *authorizeTokenRepo) FindByAccessToken(ctx context.Context, accessToken string) (*biz.AuthorizeToken, error) {
	model, err := r.data.db.AuthorizeToken.Query().
		Where(authorizetoken.AccessToken(accessToken), authorizetoken.Deleted(false)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return toBizToken(model), nil
}

// Delete soft-deletes an authorization token from the database.
func (r *authorizeTokenRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.data.db.AuthorizeToken.UpdateOneID(int(id)).
		SetDeleted(true).
		Save(ctx)
	return err
}
