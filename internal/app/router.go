package app

import (
	"github.com/labstack/echo/v4"
	"localdev.me/authorizer/config"
	"localdev.me/authorizer/internal/delivery/http/middleware"
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
	g.POST("/:id/applications/:app_id/roles", h.AssignUserRoles(), middleware.RequirePermission("AUTHORIZER-SERVICE", "user.assign_roles"))
}

func MapRolePrivateRoutes(g *echo.Group, h *v1.RoleHandler) {
	g.POST("", h.Create(), middleware.RequirePermission("AUTHORIZER-SERVICE", "role.create"))
	g.POST("/:id/permissions", h.GrantRolePermissions())
}

func MapAppPrivateRoutes(g *echo.Group, h *v1.AppHandler) {
	g.POST("", h.Create(), middleware.RequirePermission("AUTHORIZER-SERVICE", "application.create"))
}

func MapPermPrivateRoutes(g *echo.Group, h *v1.PermHandler) {
	g.POST("/sync", h.Sync(), middleware.RequirePermission("AUTHORIZER-SERVICE", "permission.sync"))
}
