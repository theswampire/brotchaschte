package server

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type GetOrdersResponse struct {
	Ok     bool    `json:"ok"`
	Orders []Order `json:"orders"`
	Error  string  `json:"error"`
}

type Order struct {
	Id             string `json:"id"`
	Lastname       string `json:"lastname"`
	Firstname      string `json:"firstname"`
	Phonenumber    string `json:"phonenumber"`
	Email          string `json:"email"`
	AmountsizeS    uint8  `json:"amountsizeS"`
	AmountsizeM    uint8  `json:"amountsizeM"`
	AmountsizeL    uint8  `json:"amountsizeL"`
	Address        string `json:"address"`
	City           string `json:"city"`
	Postalcode     string `json:"postalcode"`
	Deliverymethod string `json:"deliverymethod"`
	Remarks        string `json:"remarks"`
	Ordermethod    string `json:"ordermethod"`
	Vendorgroup    string `json:"vendorgroup"`
}

func (s *Server) getOrdersHandler(c *echo.Context) error {
	raw_orders, err := s.orders.GetOrders(c.Request().Context())
	if err != nil {
		return JSON(c, http.StatusInternalServerError, &GetOrdersResponse{Ok: false, Error: "Ein Fehler ist aufgetreten. Die Bestellungen konnten nicht geladen werden."})
	}
	orders := make([]Order, len(raw_orders))
	for i, o := range raw_orders {
		orders[i] = Order{Id: o.ID,
			Lastname:       o.Lastname,
			Firstname:      o.Firstname,
			Phonenumber:    o.Phonenumber,
			Email:          o.Email,
			AmountsizeS:    uint8(o.Amountsizes),
			AmountsizeM:    uint8(o.Amountsizem),
			AmountsizeL:    uint8(o.Amountsizel),
			Address:        o.Address,
			City:           o.City,
			Postalcode:     o.Postalcode,
			Deliverymethod: o.Deliverymethod,
			Remarks:        o.Remarks.String,
			Ordermethod:    o.Ordermethod.String,
			Vendorgroup:    o.Vendorgroup.String,
		}
	}

	return JSON(c, http.StatusOK, &GetOrdersResponse{Ok: true, Orders: orders})
}

func (s *Server) generateRandomOrdersHandler(c *echo.Context) error {
	raw_orders := s.orders.GenerateRandomOrders(c.Request().Context(), 50)
	err := s.orders.CreateBulk(c.Request().Context(), raw_orders)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, &GetOrdersResponse{Ok: false, Error: "Ein Fehler ist aufgetreten. Es konnten keine zufälligen Bestellungen generiert werden."})
	}
	orders := make([]Order, len(raw_orders))
	for i, o := range raw_orders {
		orders[i] = Order{Id: o.ID,
			Lastname:       o.Lastname,
			Firstname:      o.Firstname,
			Phonenumber:    o.Phonenumber,
			Email:          o.Email,
			AmountsizeS:    uint8(o.Amountsizes),
			AmountsizeM:    uint8(o.Amountsizem),
			AmountsizeL:    uint8(o.Amountsizel),
			Address:        o.Address,
			City:           o.City,
			Postalcode:     o.Postalcode,
			Deliverymethod: o.Deliverymethod,
			Remarks:        o.Remarks.String,
			Ordermethod:    o.Ordermethod.String,
			Vendorgroup:    o.Vendorgroup.String,
		}
	}
	return JSON(c, http.StatusOK, &GetOrdersResponse{Ok: true, Orders: orders})
}
