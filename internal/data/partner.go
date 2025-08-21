package data

import (
	"context"
	"teslatrack/internal/biz"
	"teslatrack/internal/data/ent"
	"teslatrack/internal/data/ent/partner"
)

var _ biz.PartnerRepo = (*partnerRepo)(nil)

type partnerRepo struct {
	data *Data
}

// NewPartnerRepo .
func NewPartnerRepo(data *Data) biz.PartnerRepo {
	return &partnerRepo{data: data}
}

// Get implements biz.PartnerRepo.
func (p *partnerRepo) Get(ctx context.Context, clientID string) (*biz.Partner, error) {
	po, err := p.data.db.Partner.
		Query().
		Where(partner.ClientID(clientID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return &biz.Partner{
		ID:          po.ID,
		ClientID:    po.ClientID,
		AccessToken: po.AccessToken,
		ExpiresIn:   int32(po.ExpiresIn),
		TokenType:   po.TokenType,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}, nil
}

// MustGet implements biz.PartnerRepo.
func (p *partnerRepo) MustGet(ctx context.Context, clientID string) (*biz.Partner, error) {
	partner, err := p.Get(ctx, clientID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return partner, err
}

// Create implements biz.PartnerRepo.
func (p *partnerRepo) Create(ctx context.Context, b *biz.Partner) error {
	_, err := p.data.db.Partner.
		Create().
		SetClientID(b.ClientID).
		SetAccessToken(b.AccessToken).
		SetExpiresIn(int(b.ExpiresIn)).
		SetTokenType(b.TokenType).
		Save(ctx)
	return err
}

// Update implements biz.PartnerRepo.
func (p *partnerRepo) Update(ctx context.Context, id int, b *biz.Partner) error {
	_, err := p.data.db.Partner.
		UpdateOneID(id).
		SetAccessToken(b.AccessToken).
		SetExpiresIn(int(b.ExpiresIn)).
		SetTokenType(b.TokenType).
		Save(ctx)
	return err
}
