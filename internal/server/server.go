package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"brotchaschte2/internal/database"
	"brotchaschte2/internal/stores/auth"
	"brotchaschte2/internal/stores/delivery"
	"brotchaschte2/internal/stores/geocode"
	"brotchaschte2/internal/stores/orders"
	"brotchaschte2/internal/stores/pdfexport"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int

	db       database.Service
	auth     auth.AuthService
	orders   orders.OrderStore
	geocode  geocode.GeocodeStore
	delivery delivery.DeliveryStore
	pdf      pdfexport.PDFExportStore
}

func NewServer() *http.Server {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
	}

	NewServer := &Server{
		port: port,

		db:       database.New(),
		auth:     auth.New(),
		orders:   orders.New(),
		geocode:  geocode.New(),
		delivery: delivery.New(),
		pdf:      pdfexport.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server

}
