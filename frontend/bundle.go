//go:build bundle

package frontend

import (
	"embed"
	"log"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

//go:embed "dist"
var bundle embed.FS

func RegisterFrontend(e *echo.Echo) {
	log.Println("registering frontend routes")
	e.Filesystem = echo.MustSubFS(bundle, "dist")

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   ".",
		Browse: false,
		HTML5:  true,
	}))
}
