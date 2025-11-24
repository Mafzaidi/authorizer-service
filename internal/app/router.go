package app

import (
	"github.com/labstack/echo/v4"
	"localdev.me/authorizer/config"
	v1 "localdev.me/authorizer/internal/delivery/http/v1"
)

func MapAuthPublicRoutes(g *echo.Group, h *v1.AuthHandler) {
	g.POST("/login", h.Login())
}

func MapAuthPrivateRoutes(g *echo.Group, h *v1.AuthHandler) {
	g.POST("/logout", h.Logout())
}

func MapUserPublicRoutes(g *echo.Group, h *v1.UserHandler, cfg *config.Config) {
	g.POST("", h.RegisterUser())
}

func MapUserPrivateRoutes(g *echo.Group, h *v1.UserHandler) {
	g.GET("/:id", h.GetUserProfile())
	g.GET("", h.GetUserList())
	g.PATCH("/:id", h.UpdateUserProfile())
}

func MapRolePrivateRoutes(g *echo.Group, h *v1.RoleHandler) {
	g.POST("", h.Create())
}
