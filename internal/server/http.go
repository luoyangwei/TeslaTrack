package server

import (
	v1 "teslatrack/api/teslatrack/v1"
	"teslatrack/internal/biz"
	"teslatrack/internal/conf"
	"teslatrack/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	c *conf.Server,
	logger log.Logger,
	redirector *Redirector,
	partnerUsecase *biz.PartnerUsecase,
	authorize *service.AuthorizeService,
) *kratoshttp.Server {
	// Define server options.
	var opts = []kratoshttp.ServerOption{
		// Add middleware for recovery and logging.
		kratoshttp.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
		),
		// Add a filter for redirection.
		kratoshttp.Filter(redirector.RedirectFilter),
	}
	// Set network, address and timeout from config.
	if c.Http.Network != "" {
		opts = append(opts, kratoshttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, kratoshttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, kratoshttp.Timeout(c.Http.Timeout.AsDuration()))
	}

	// Create a new HTTP server.
	srv := kratoshttp.NewServer(opts...)

	// Register static file server.
	// srv.HandlePrefix("/", NewStaticServer())

	// Register the Authorize service.
	v1.RegisterAuthorizeHTTPServer(srv, authorize)

	// Initialize partner usecase.
	if err := partnerUsecase.Initialize(); err != nil {
		panic(err)
	}

	return srv
}
