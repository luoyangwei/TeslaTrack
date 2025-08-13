package service

import (
	"context"

	v1 "teslatrack/api/teslatrack/v1"
	"teslatrack/internal/biz"
)

type AuthorizeService struct {
	v1.UnimplementedAuthorizeServer

	uc *biz.AuthorizeUsecase
}

func NewAuthorizeService(uc *biz.AuthorizeUsecase) *AuthorizeService {
	return &AuthorizeService{uc: uc}
}

func (s *AuthorizeService) CreateAuthorize(ctx context.Context, req *v1.CreateAuthorizeRequest) (*v1.CreateAuthorizeReply, error) {
	err := s.uc.Create(ctx, &biz.Authorize{
		ClientID:     req.ClientId,
		ClientSecret: req.ClientSecret,
		GrantType:    req.GrantType,
		RedirectURI:  req.RedirectURI,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateAuthorizeReply{}, nil
}
