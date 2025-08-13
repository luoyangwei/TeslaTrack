package data

import (
	"context"
	"teslatrack/internal/biz"
	"teslatrack/internal/data/ent/authorize"
)

// A compile-time check to ensure that authorizeRepo implements the biz.AuthorizeRepo interface.
var _ biz.AuthorizeRepo = (*authorizeRepo)(nil)

// authorizeRepo is the data access layer implementation for authorization.
// It uses the ent framework to interact with the database.
type authorizeRepo struct {
	data *Data // data field holds the database client and other data sources.
}

// NewAuthorizeRepo creates a new authorizeRepo.
// It takes a Data struct which contains the database client.
func NewAuthorizeRepo(data *Data) biz.AuthorizeRepo {
	return &authorizeRepo{data}
}

// Create saves a new authorization record to the database.
// It maps the biz.Authorize model to an ent.Authorize create operation.
func (repo *authorizeRepo) Create(ctx context.Context, auth *biz.Authorize) error {
	_, err := repo.data.db.Authorize.Create().
		SetClientID(auth.ClientID).
		SetClientSecret(auth.ClientSecret).
		SetGrantType(auth.GrantType).
		SetRedirectURI(auth.RedirectURI).
		Save(ctx)
	return err
}

// FindByClientID retrieves an authorization record from the database by its client ID.
// It queries the database and maps the resulting ent.Authorize model to a biz.Authorize model.
func (repo *authorizeRepo) FindByClientID(ctx context.Context, clientID string) (*biz.Authorize, error) {
	model, err := repo.data.db.Authorize.Query().
		Where(authorize.ClientID(clientID)).
		First(ctx)
	if err != nil {
		return nil, err // Return error if the query fails or no record is found.
	}
	// Map the ent model to the biz model.
	return &biz.Authorize{
		ID:           int64(model.ID),
		ClientID:     model.ClientID,
		ClientSecret: model.ClientSecret,
		GrantType:    model.GrantType,
		RedirectURI:  model.RedirectURI,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}, nil
}

// Update modifies an existing authorization record in the database.
// It finds the record by its ID and updates the specified fields.
func (repo *authorizeRepo) Update(ctx context.Context, auth *biz.Authorize) error {
	_, err := repo.data.db.Authorize.UpdateOneID(int(auth.ID)).
		// The original comment `// 不能修改` means "cannot be modified".
		// ClientID and ClientSecret are generally considered immutable and should not be updated.
		// SetClientID(auth.ClientID).
		// SetClientSecret(auth.ClientSecret).
		SetGrantType(auth.GrantType).
		SetRedirectURI(auth.RedirectURI).
		Save(ctx)
	return err
}
