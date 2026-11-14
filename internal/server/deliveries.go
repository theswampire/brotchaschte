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

func (s *Server) getDeliveriesHandler(c *echo.Context) error {
	deliveries_row, err := s.delivery.GetDeliveries(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &GetDeliveriesRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Die Liefer-Touren konnten nicht geladen werden."})
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
	return c.JSON(http.StatusOK, &GetDeliveriesRes{Ok: true, Deliveries: deliveries})
}

type CreateDeliveryRes struct {
	Ok       bool        `json:"ok"`
	Error    string      `json:"error"`
	Delivery DeliveryDTO `json:"delivery"`
}

func (s *Server) createDeliveryHandler(c *echo.Context) error {
	until := c.FormValue("until")
	color := c.FormValue("color")
	label := c.FormValue("label")

	delivery, err := s.delivery.CreateNewDelivery(c.Request().Context(), label, until, color)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &CreateDeliveryRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Es konnte keine neue Liefertour erstellt werden."})
	}
	return c.JSON(http.StatusOK, &CreateDeliveryRes{Ok: true,
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
	Orders   []OrderDTO  `json:"orders"`
}

func (s *Server) getDeliveryHandler(c *echo.Context) error {
	param := c.QueryParam("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &GetDeliveryRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Falscher Datentyp."})
	}
	delivery, err := s.delivery.GetDelivery(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &GetDeliveryRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Die Lieferung konnte nicht geladen werden."})
	}
	raw_orders, err := s.delivery.GetOrdersOfDelivery(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &GetDeliveryRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Die Bestellungen für die Lieferung konnten nicht geladen werden."})
	}

	orders := make([]OrderDTO, len(raw_orders))
	for i, order := range raw_orders {
		orders[i] = OrderDTO{
			Id:             order.ID,
			Lastname:       order.Lastname,
			Firstname:      order.Firstname,
			Phonenumber:    order.Phonenumber,
			Email:          order.Email,
			AmountsizeS:    uint8(order.Amountsizes),
			AmountsizeM:    uint8(order.Amountsizem),
			AmountsizeL:    uint8(order.Amountsizel),
			Address:        order.Address,
			City:           order.City,
			Postalcode:     order.Postalcode,
			Deliverymethod: order.Deliverymethod,
			Remarks:        order.Remarks.String,
			Ordermethod:    order.Ordermethod.String,
			Vendorgroup:    order.Vendorgroup.String,
		}
	}

	return c.JSON(http.StatusOK, &GetDeliveryRes{Ok: true, Orders: orders, Delivery: DeliveryDTO{
		ID:    delivery.ID,
		Label: delivery.Label,
		Until: delivery.Until,
		Color: delivery.Color,
	}})
}

type AssignOrdersRes struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type AssignOrdersBody struct {
	DeliveryId         int64    `json:"delivery_id"`
	LenSelected        int      `json:"len_selected"`
	LenDeselected      int      `json:"len_deselected"`
	SelectedOrderIds   []string `json:"selected_order_ids"`
	DeselectedOrderIds []string `json:"deselected_order_ids"`
}

func (s *Server) assignOrdersHandler(c *echo.Context) error {
	var data *AssignOrdersBody
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, &AssignOrdersRes{Ok: false, Error: "Ein Fehler ist aufgetreten."})
	}
	if data.SelectedOrderIds == nil || data.DeselectedOrderIds == nil {
		return c.JSON(http.StatusBadRequest, &AssignOrdersRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Listen Null."})
	}
	if len(data.SelectedOrderIds) != data.LenSelected {
		return c.JSON(http.StatusBadRequest, &AssignOrdersRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Die Länge der selected unstimmig."})
	}
	if len(data.DeselectedOrderIds) != data.LenDeselected {
		return c.JSON(http.StatusBadRequest, &AssignOrdersRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Die Länge der deselected unstimmig."})
	}

	if err := s.delivery.AssignOrders(c.Request().Context(), data.DeliveryId, data.SelectedOrderIds, data.DeselectedOrderIds); err != nil {
		return c.JSON(http.StatusInternalServerError, &AssignOrdersRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Die Bestellungen konnten nicht der Lieferung zugewiesen werden."})
	}
	return c.JSON(http.StatusOK, &AssignOrdersRes{Ok: true})
}

type EditDeliveryRes struct {
	Ok       bool        `json:"ok"`
	Error    string      `json:"error"`
	Delivery DeliveryDTO `json:"delivery"`
}

func (s *Server) editDeliveryHandler(c *echo.Context) error {
	id_str := c.FormValue("id")
	until := c.FormValue("until")
	color := c.FormValue("color")
	label := c.FormValue("label")
	id, err := strconv.ParseInt(id_str, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &EditDeliveryRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Eine ungültige Liefertour-ID."})
	}
	if until == "" || color == "" {
		return c.JSON(http.StatusBadRequest, &EditDeliveryRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Eine ungültige Farbe oder Zeit."})
	}
	delivery, err := s.delivery.UpdateDelivery(c.Request().Context(), id, until, label, color)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &EditDeliveryRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Die Lieferung konnte nicht verändert werden."})
	}

	return c.JSON(http.StatusOK, &EditDeliveryRes{Ok: true, Delivery: DeliveryDTO{ID: delivery.ID, Color: delivery.Color, Label: delivery.Label, Until: delivery.Until}})
}
