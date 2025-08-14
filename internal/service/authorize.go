package service

import (
	"context"

	v1 "teslatrack/api/teslatrack/v1"
	"teslatrack/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// AuthorizeService is the service implementation for the Authorize gRPC API.
// It orchestrates the authorization flow by calling the appropriate business logic use cases.
type AuthorizeService struct {
	v1.UnimplementedAuthorizeServer

	uc  *biz.AuthorizeUsecase
	log *log.Helper
}

// NewAuthorizeService creates a new AuthorizeService.
func NewAuthorizeService(uc *biz.AuthorizeUsecase, logger log.Logger) *AuthorizeService {
	return &AuthorizeService{uc: uc, log: log.NewHelper(logger)}
}

// CreateAuthorize handles the RPC for creating a new client configuration.
func (s *AuthorizeService) CreateAuthorize(ctx context.Context, req *v1.CreateAuthorizeRequest) (*v1.CreateAuthorizeReply, error) {
	// Map the gRPC request to the business layer's Authorize model.
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

// Redirect handles the RPC for initiating the authorization redirect.
// It calls the Redirect use case and maps the result to the gRPC reply.
func (s *AuthorizeService) Redirect(ctx context.Context, req *v1.RedirectRequest) (*v1.RedirectReply, error) {
	// Call the business logic to get redirect parameters.
	redirect, err := s.uc.Redirect(ctx, req.ClientId)
	if err != nil {
		return nil, err
	}

	// Map the business object to the gRPC reply message.
	return &v1.RedirectReply{
		Scope:                  redirect.Scope,
		State:                  redirect.State,
		Nonce:                  redirect.Nonce,
		PromptMissingScopes:    redirect.PromptMissingScopes,
		RequireRequestedScopes: redirect.RequireRequestedScopes,
		RedirectUri:            redirect.RedirectURI,
	}, nil
}

// Callback handles the RPC for the OAuth 2.0 callback.
// It receives the authorization code from the client.
func (s *AuthorizeService) Callback(ctx context.Context, req *v1.CallbackRequest) (*v1.CallbackReply, error) {
	s.log.Infow("msg", "Tesla callback code", "code", req.Code)
	// Call the business logic to handle the authorization code.
	err := s.uc.Callback(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return &v1.CallbackReply{}, nil
}
