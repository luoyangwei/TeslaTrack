package server

import (
	v1 "teslatrack/api/helloworld/v1"
	"teslatrack/internal/conf"
	"teslatrack/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
)

// NewHTTPServer creates a new HTTP server with gRPC-Gateway integration.
// This server handles both the main API endpoints and integrates with a Mux router
// for additional custom routes.
//
// Parameters:
//   - c: Server configuration containing HTTP-specific settings (network, address, timeout)
//   - greeter: Greeter service implementation for handling API requests
//   - muxRouter: Gorilla Mux router for custom route handling
//   - logger: Logger instance for server operations
//
// Returns: *http.Server - A configured HTTP server with gRPC-Gateway and Mux integration
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, muxRouter *mux.Router, logger log.Logger) *http.Server {
	// Configure server options with recovery middleware
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}

	// Apply network configuration if specified
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}

	// Apply address configuration if specified
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}

	// Apply timeout configuration if specified
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	// Create the HTTP server with configured options
	srv := http.NewServer(opts...)

	// Mount the Mux router under the /mux path for custom routes
	srv.HandlePrefix("/", muxRouter)

	// Register the gRPC-Gateway HTTP server with the greeter service
	v1.RegisterGreeterHTTPServer(srv, greeter)

	return srv
}
