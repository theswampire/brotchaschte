package server

import (
	"brotchaschte2/frontend"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(s.loadAndSaveHttpSession)
	e.Use(s.addAuthCookie)
	e.GET("/api", s.requireSession(func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	}))
	s.registerApiV1Routes(e)

	frontend.RegisterFrontend(e)
	return e

}

func (s *Server) registerApiV1Routes(e *echo.Echo) {
	router := e.Group("/api/v1")

	auth := router.Group("/auth")
	auth.POST("/login", s.loginHandler)
	auth.Any("/logout", s.logoutHandler)
	auth.Any("/check", s.checkAuthenticatedHandler)

	kek := router.Group("/kek")
	kek.Use(s.requireSession)
	kek.Any("", s.kekHandler)

	orders := router.Group("/orders")
	orders.Use(s.requireSession)
	orders.GET("", s.getOrdersHandler)
	orders.POST("/generate", s.generateRandomOrdersHandler)
	orders.POST("/geolocate", s.geolocateOrdersHandler)
	orders.POST("/sync", s.pullFromUpstreamHandler)
	orders.GET("/sorted", s.getDeliveryOrdersSorted)

	deliveries := router.Group("/deliveries")
	deliveries.Use(s.requireSession)
	deliveries.GET("", s.getDeliveryHandler)
	deliveries.POST("", s.createDeliveryHandler)
	deliveries.PATCH("", s.editDeliveryHandler)
	deliveries.GET("/all", s.getDeliveriesHandler)
	deliveries.POST("/assign", s.assignOrdersHandler)
}
