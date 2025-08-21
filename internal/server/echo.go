package server

import (
	"github.com/labstack/echo/v4"
)

// NewStaticServer creates a new echo server for serving static files from the "web" directory.
func NewStaticServer() *echo.Echo {
	router := echo.New()
	router.Static("/", "web")
	return router
}
