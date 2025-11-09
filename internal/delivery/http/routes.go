package http

import (
	"github.com/labstack/echo/v4"
	"localdev.me/authorizer/config"
)

func MapPublicRoutes(g *echo.Group, h User, cfg *config.Config) {
	g.POST("/register", h.Register())
}
