package server

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type GeolocateResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

func (s *Server) geolocateOrdersHandler(c *echo.Context) error {
	err := s.geocode.GeocodeOrders(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusOK, &GeolocateResponse{Ok: false, Error: "Ein Fehler ist aufgetreten. Die Bestellungen konnten nicht geokodiert werden."})
	}

	return c.JSON(http.StatusOK, &GeolocateResponse{Ok: true})
}
