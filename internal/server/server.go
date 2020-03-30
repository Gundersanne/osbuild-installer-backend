package Server

import (
	//	"fmt"
	"net/http"

	"github.com/osbuild/osbuild-installer-backend/internal/api"

	"github.com/labstack/echo/v4"
)


type Handlers struct { }

func (s *Handlers) OsbuildInstallerViewsV1GetVersion(c echo.Context) error {
	return c.String(http.StatusOK, "1")
}

func Run(address string, routePrefix string) {
	echo := echo.New()
	routeGroup := echo.Group(routePrefix)

	var s Handlers
	Api.RegisterHandlers(routeGroup, &s) //  <<<<<-
	echo.Start(address)
}
