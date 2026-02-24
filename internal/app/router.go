package app

import (
	"github.com/labstack/echo/v4"
	"github.com/mafzaidi/authorizer/config"
	"github.com/mafzaidi/authorizer/internal/delivery/http/middleware"
	v1 "github.com/mafzaidi/authorizer/internal/delivery/http/v1"
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
	g.POST("/:id/applications/:app_id/roles", h.AssignUserRoles(), middleware.RequirePermission("AUTHORIZER", "user.assign_roles"))
}

func MapRolePrivateRoutes(g *echo.Group, h *v1.RoleHandler) {
	g.POST("", h.Create(), middleware.RequirePermission("AUTHORIZER", "role.create"))
	g.POST("/:id/permissions", h.GrantRolePermissions())
}

func MapAppPrivateRoutes(g *echo.Group, h *v1.AppHandler) {
	g.POST("", h.Create(), middleware.RequirePermission("AUTHORIZER", "application.create"))
}

func MapPermPrivateRoutes(g *echo.Group, h *v1.PermHandler) {
	g.POST("/sync", h.Sync(), middleware.RequirePermission("AUTHORIZER", "permission.sync"))
}
