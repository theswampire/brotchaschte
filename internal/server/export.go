package server

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
)

func (s *Server) exportDeliveriesToPDFHandler(c *echo.Context) error {
	download := c.QueryParam("download")
	w := c.Response()
	w.Header().Set("Content-Type", "application/pdf")
	if download == "true" {
		w.Header().Set("Content-Disposition", "attachment; filename=\"lieferungen.pdf\"")
	}
	err := s.pdf.DeliveryRuns(c.Request().Context(), w)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Ein Fehler ist aufgetreten.")
	}
	return nil
}
