package server

import (
	"brotchaschte2/frontend"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.GET("/api", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	frontend.RegisterFrontend(e)
	return e

}
