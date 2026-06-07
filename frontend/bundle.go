//go:build bundle

package frontend

import (
	"embed"
	"log"

	"github.com/labstack/echo/v5"
)

//go:embed "dist"
var bundle embed.FS

func RegisterFrontend(e *echo.Echo) {
	log.Println("registering frontend routes")
	e.Filesystem = echo.MustSubFS(bundle, "dist")

	e.Static("/", ".")
}
