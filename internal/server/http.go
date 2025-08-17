package server

import (
	v1 "teslatrack/api/teslatrack/v1"
	"teslatrack/internal/conf"
	"teslatrack/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, redirector *Redirector, authorize *service.AuthorizeService, logger log.Logger) *kratoshttp.Server {
	var opts = []kratoshttp.ServerOption{
		kratoshttp.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
		),
		kratoshttp.Filter(redirector.RedirectFilter),
	}
	if c.Http.Network != "" {
		opts = append(opts, kratoshttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, kratoshttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, kratoshttp.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := kratoshttp.NewServer(opts...)
	v1.RegisterAuthorizeHTTPServer(srv, authorize)

	return srv
}
