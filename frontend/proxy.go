//go:build !bundle

package frontend

import (
	"net/url"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func RegisterFrontend(e *echo.Echo) {
	upstream, _ := url.Parse("http://localhost:5173")
	proxy := middleware.ProxyWithConfig(middleware.ProxyConfig{
		Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
			{URL: upstream},
		}),
	})

	e.Any("/*", func(c *echo.Context) error {
		return proxy(func(c *echo.Context) error { return nil })(c)
	})
}
