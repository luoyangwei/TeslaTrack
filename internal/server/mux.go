package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// NewMuxRouter creates a new Gorilla Mux router with predefined routes.
// This router is designed to be integrated into the main HTTP server.
// Returns: *mux.Router - A configured router with home route
func NewMuxRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/home", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "Hello Gorilla Mux!")
	}).Methods("GET")
	return router
}

// NewMuxServer creates a standalone HTTP server using Gorilla Mux.
// This server runs independently from the main HTTP server and can be used
// for testing or as a separate service endpoint.
//
// Parameters:
//   - c: Server configuration containing network and timeout settings
//   - logger: Logger instance for server operations
//
// Returns: *khttp.Server - A configured HTTP server with Mux router
// func NewMuxServer(c *conf.Server, logger log.Logger) *khttp.Server {
// 	// Create a new Mux router for this standalone server
// 	router := mux.NewRouter()
// 	router.HandleFunc("/home", func(w http.ResponseWriter, req *http.Request) {
// 		fmt.Fprint(w, "Hello Gorilla Mux!")
// 	}).Methods("GET")

// 	// Configure server options with recovery middleware
// 	var opts = []khttp.ServerOption{
// 		khttp.Middleware(
// 			recovery.Recovery(),
// 		),
// 		khttp.Address(":8001"), // Use a different port to avoid conflicts
// 	}

// 	// Create and configure the HTTP server
// 	httpSrv := khttp.NewServer(opts...)
// 	httpSrv.HandlePrefix("/", router)
// 	return httpSrv
// }
