package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
)

type MaybeFloat64 struct {
	Ok    bool    `json:"ok"`
	Value float64 `json:"value"`
}

func NewMaybeFloat64(v sql.NullFloat64) MaybeFloat64 {
	return MaybeFloat64{
		Ok:    v.Valid,
		Value: v.Float64,
	}
}

type GetOrdersResponse struct {
	Ok                   bool                     `json:"ok"`
	Error                string                   `json:"error"`
	NotGeocodedOrders    []OrderDTO               `json:"not_geocoded_orders"`
	FailedGeocodedOrders []FailedGeocodedOrderDTO `json:"failed_geocoded_orders"`
	GeocodedOrders       []GeocodedOrderDTO       `json:"geocoded_orders"`
}

type OrderDTO struct {
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

type GeocodedOrderDTO struct {
	OrderDTO

	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Label string  `json:"label"`
}

type FailedGeocodedOrderDTO struct {
	OrderDTO

	Message           string       `json:"message"`
	SuggestedNorthing MaybeFloat64 `json:"suggested_northing"`
	SuggestedEasting  MaybeFloat64 `json:"suggested_easting"`
	SuggestedLabel    string       `json:"suggested_label"`
}

type OrderDeliveryDTO struct {
	OrderDTO
	X                     float64 `json:"x"`
	Y                     float64 `json:"y"`
	Label                 string  `json:"label"`
	HasAssociatedDelivery bool    `json:"has_associated_delivery"`
	DeliveryId            int64   `json:"delivery_id"`
	DeliveryLabel         string  `json:"delivery_label"`
	DeliveryColor         string  `json:"delivery_color"`
}

func (s *Server) getOrdersHandler(c *echo.Context) error {
	not_geocoded_rows, err := s.geocode.GetNotGeocodedExcludeFailed(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &GetOrdersResponse{Ok: false, Error: "Ein Fehler ist aufgetreten. Die Bestellungen konnten nicht geladen werden."})
	}
	failed_geocoded_rows, err := s.geocode.GetFailedGeocodeOrders(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &GetOrdersResponse{Ok: false, Error: "Ein Fehler ist aufgetreten. Die Bestellungen mit fehlgeschlagener Geolokalisierung konnten nicht geladen werden."})
	}
	geocoded_rows, err := s.geocode.GetGeocodedOrders(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &GetOrdersResponse{Ok: false, Error: "Ein Fehler ist aufgetreten. Die Bestellungen mit Geolokalisierung konnten nicht geladen werden."})
	}

	not_geocoded := make([]OrderDTO, len(not_geocoded_rows))
	for i, o := range not_geocoded_rows {
		not_geocoded[i] = OrderDTO{Id: o.ID,
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
	failed_geocodes := make([]FailedGeocodedOrderDTO, len(failed_geocoded_rows))
	for i, o := range failed_geocoded_rows {
		failed_geocodes[i] = FailedGeocodedOrderDTO{
			OrderDTO: OrderDTO{Id: o.ID,
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
			},
			Message:           o.Message,
			SuggestedNorthing: NewMaybeFloat64(o.SuggestedNorthing),
			SuggestedEasting:  NewMaybeFloat64(o.SuggestedEasting),
			SuggestedLabel:    o.SuggestedLabel.String,
		}
	}
	geocoded_orders := make([]GeocodedOrderDTO, len(geocoded_rows))
	for i, o := range geocoded_rows {
		geocoded_orders[i] = GeocodedOrderDTO{
			OrderDTO: OrderDTO{Id: o.ID,
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
			},
			X:     o.X,
			Y:     o.Y,
			Label: o.Label,
		}
	}
	return c.JSON(http.StatusOK, &GetOrdersResponse{Ok: true, NotGeocodedOrders: not_geocoded, FailedGeocodedOrders: failed_geocodes, GeocodedOrders: geocoded_orders})
}

func (s *Server) generateRandomOrdersHandler(c *echo.Context) error {
	raw_orders := s.orders.GenerateRandomOrders(c.Request().Context(), 50)
	err := s.orders.CreateBulk(c.Request().Context(), raw_orders)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &GetOrdersResponse{Ok: false, Error: "Ein Fehler ist aufgetreten. Es konnten keine zufälligen Bestellungen generiert werden."})
	}
	orders := make([]OrderDTO, len(raw_orders))
	for i, o := range raw_orders {
		orders[i] = OrderDTO{Id: o.ID,
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
	return c.JSON(http.StatusOK, &GetOrdersResponse{Ok: true, NotGeocodedOrders: orders})
}

func (s *Server) pullFromUpstreamHandler(c *echo.Context) error {
	err := s.orders.SyncWithUpstream(c.Request().Context())
	if err != nil {
		log.Printf("error: failed to sync with upstream: %s\n", err.Error())
		return c.JSON(http.StatusInternalServerError, &GetOrdersResponse{Ok: false, Error: "Ein Fehler ist aufgetreten. Bestellungen konnten nicht importiert werden."})
	}

	return c.JSON(http.StatusInternalServerError, &GetOrdersResponse{Ok: true})
}

type GetDeliveryOrdersRes struct {
	Ok         bool               `json:"ok"`
	Error      string             `json:"error"`
	Orders     []OrderDeliveryDTO `json:"orders"`
	Deliveries []DeliveryDTO      `json:"deliveries"`
}

func (s *Server) getDeliveryOrdersSorted(c *echo.Context) error {
	orders_row, err := s.delivery.GetOrderCoordinates(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &GetDeliveryOrdersRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Die sortierten Bestellungen konnten nicht geladen werden."})
	}
	delivery_row, err := s.delivery.GetDeliveries(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &GetDeliveryOrdersRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Die Lieferungen zu den sortierten Bestellungen konnten nicht geladen werden."})
	}
	orders := make([]OrderDeliveryDTO, len(orders_row))
	for i, o := range orders_row {
		orders[i] = OrderDeliveryDTO{
			OrderDTO: OrderDTO{
				Id:             o.ID,
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
			},
			X:                     o.X,
			Y:                     o.Y,
			Label:                 o.AddressLabel,
			HasAssociatedDelivery: o.DeliveryID.Valid,
			DeliveryId:            o.DeliveryID.Int64,
			DeliveryLabel:         o.DeliveryLabel.String,
			DeliveryColor:         o.Color.String,
		}
	}
	deliveries := make([]DeliveryDTO, len(delivery_row))
	for i, d := range delivery_row {
		deliveries[i] = DeliveryDTO{ID: d.ID, Label: d.Label, Until: d.Until, Color: d.Color}
	}
	return c.JSON(http.StatusOK, &GetDeliveryOrdersRes{Ok: true, Orders: orders, Deliveries: deliveries})
}
