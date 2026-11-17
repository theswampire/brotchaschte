package server

import (
	"log"
	"net/http"
	"strconv"

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

type ManualLocateRes struct {
	Ok    bool
	Error string
}

func (s *Server) manualLocateHandler(c *echo.Context) error {
	order_id := c.FormValue("order_id")
	northing_str := c.FormValue("northing")
	easting_str := c.FormValue("easting")
	northing, err := strconv.ParseFloat(northing_str, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ManualLocateRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Ungültige Koordinate."})
	}
	easting, err := strconv.ParseFloat(easting_str, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ManualLocateRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Ungültige Koordinate."})
	}
	_, err = s.geocode.ManuallyLocateOrder(c.Request().Context(), order_id, northing, easting)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ManualLocateRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Die Koordinaten konnten nicht gespeichert werden."})
	}

	return c.JSON(http.StatusOK, &ManualLocateRes{Ok: true})
}

func (s *Server) manualLocateWithAddressUpdateHandler(c *echo.Context) error {
	order_id := c.FormValue("order_id")
	northing_str := c.FormValue("northing")
	easting_str := c.FormValue("easting")
	address := c.FormValue("address")
	postalcode := c.FormValue("postalcode")
	city := c.FormValue("city")
	remarks := c.FormValue("remarks")

	northing, err := strconv.ParseFloat(northing_str, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ManualLocateRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Ungültige Koordinate."})
	}
	easting, err := strconv.ParseFloat(easting_str, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ManualLocateRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Ungültige Koordinate."})
	}
	_, err = s.geocode.ManuallyLocateAndUpdateAddressOrder(c.Request().Context(), order_id, address, postalcode, city, remarks, northing, easting)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, &ManualLocateRes{Ok: false, Error: "Ein Fehler ist aufgetreten. Die Koordinaten konnten nicht gespeichert werden."})
	}

	return c.JSON(http.StatusOK, &ManualLocateRes{Ok: true})
}
