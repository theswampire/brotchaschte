package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
)

func JSON(c *echo.Context, code int, i any) error {
	m, err := json.Marshal(i)
	if err != nil {
		log.Println("JSON Serialize failed")
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.String(code, string(m))
}
