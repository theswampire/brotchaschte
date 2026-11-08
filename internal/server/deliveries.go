package server

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type DeliveriesDTO struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
	Until string `json:"until"`
	Color string `json:"color"`
}
type GetDeliveriesRes struct {
	Ok         bool            `json:"ok"`
	Error      string          `json:"error"`
	Deliveries []DeliveriesDTO `json:"deliveries"`
}

func (s *Server) getDeliveries(c *echo.Context) error {
	deliveries_row, err := s.delivery.GetDeliveries(c.Request().Context())
	if err != nil {
		return JSON(c, http.StatusInternalServerError, &GetDeliveriesRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Die Liefer-Touren konnten nicht geladen werden."})
	}
	deliveries := make([]DeliveriesDTO, len(deliveries_row))
	for i, d := range deliveries_row {
		deliveries[i] = DeliveriesDTO{
			ID:    d.ID,
			Label: d.Label,
			Until: d.Until,
			Color: d.Color,
		}
	}
	return JSON(c, http.StatusOK, &GetDeliveriesRes{Ok: true, Deliveries: deliveries})
}
