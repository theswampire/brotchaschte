package server

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (s *Server) kekHandler(c *echo.Context) error {
	return c.String(http.StatusOK, "kekekekek")
}
