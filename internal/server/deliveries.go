package server

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

type DeliveryDTO struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
	Until string `json:"until"`
	Color string `json:"color"`
}
type GetDeliveriesRes struct {
	Ok         bool          `json:"ok"`
	Error      string        `json:"error"`
	Deliveries []DeliveryDTO `json:"deliveries"`
}

func (s *Server) getDeliveries(c *echo.Context) error {
	deliveries_row, err := s.delivery.GetDeliveries(c.Request().Context())
	if err != nil {
		return JSON(c, http.StatusInternalServerError, &GetDeliveriesRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Die Liefer-Touren konnten nicht geladen werden."})
	}
	deliveries := make([]DeliveryDTO, len(deliveries_row))
	for i, d := range deliveries_row {
		deliveries[i] = DeliveryDTO{
			ID:    d.ID,
			Label: d.Label,
			Until: d.Until,
			Color: d.Color,
		}
	}
	return JSON(c, http.StatusOK, &GetDeliveriesRes{Ok: true, Deliveries: deliveries})
}

type CreateDeliveryRes struct {
	Ok       bool        `json:"ok"`
	Error    string      `json:"error"`
	Delivery DeliveryDTO `json:"delivery"`
}

func (s *Server) createDelivery(c *echo.Context) error {
	until := c.FormValue("until")
	color := c.FormValue("color")
	label := c.FormValue("label")

	delivery, err := s.delivery.CreateNewDelivery(c.Request().Context(), label, until, color)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, &CreateDeliveryRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Es konnte keine neue Liefertour erstellt werden."})
	}
	return JSON(c, http.StatusOK, &CreateDeliveryRes{Ok: true,
		Delivery: DeliveryDTO{
			ID:    delivery.ID,
			Label: delivery.Label,
			Until: delivery.Until,
			Color: delivery.Color,
		}})
}

type GetDeliveryRes struct {
	Ok       bool        `json:"ok"`
	Error    string      `json:"error"`
	Delivery DeliveryDTO `json:"delivery"`
}

func (s *Server) getDelivery(c *echo.Context) error {
	param := c.QueryParam("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return JSON(c, http.StatusBadRequest, &GetDeliveryRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Falscher Datentyp."})
	}
	delivery, err := s.delivery.GetDelivery(c.Request().Context(), id)
	if err != nil {
		return JSON(c, http.StatusBadRequest, &GetDeliveryRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Die Lieferung konnte nicht geladen werden."})
	}
	return JSON(c, http.StatusOK, &GetDeliveryRes{Ok: true, Delivery: DeliveryDTO{
		ID:    delivery.ID,
		Label: delivery.Label,
		Until: delivery.Until,
		Color: delivery.Color,
	}})
}
