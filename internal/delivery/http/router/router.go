package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mafzaidi/authorizer/internal/delivery/http/handler"
	appMiddleware "github.com/mafzaidi/authorizer/internal/delivery/http/middleware"
	"github.com/mafzaidi/authorizer/internal/infrastructure/logger"
)

// RouterConfig holds all dependencies needed for setting up routes
type RouterConfig struct {
	// Handlers
	AuthHandler   *handler.AuthHandler
	UserHandler   *handler.UserHandler
	RoleHandler   *handler.RoleHandler
	PermHandler   *handler.PermHandler
	AppHandler    *handler.AppHandler
	HealthHandler *handler.HealthHandler

	// Middleware
	JWTMiddleware echo.MiddlewareFunc

	// Logger
	Logger *logger.Logger
}

// Setup configures all routes with handlers and middleware
func Setup(e *echo.Echo, cfg *RouterConfig) error {
	// Setup CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
		AllowMethods:     []string{echo.POST, echo.GET, echo.OPTIONS, echo.PATCH, echo.PUT, echo.DELETE},
	}))

	// Setup basic middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// API version group
	v1 := e.Group("authorizer/v1")

	// Public routes group
	public := v1.Group("")

	// Public auth routes
	pblAuth := public.Group("/auth")
	mapAuthPublicRoutes(pblAuth, cfg.AuthHandler)

	// Public user routes
	pblUser := public.Group("/users")
	mapUserPublicRoutes(pblUser, cfg.UserHandler)

	// Public health routes
	pblHealth := public.Group("/health")
	pblHealth.GET("", cfg.HealthHandler.Check())

	// JWKS endpoint (public, outside of /v1)
	e.GET("/.well-known/jwks.json", cfg.AuthHandler.GetJWKS())

	// Private routes group (with JWT middleware)
	private := v1.Group("")
	private.Use(cfg.JWTMiddleware)

	// Private auth routes
	pvtAuth := private.Group("/auth")
	mapAuthPrivateRoutes(pvtAuth, cfg.AuthHandler)

	// Private user routes
	pvtUser := private.Group("/users")
	mapUserPrivateRoutes(pvtUser, cfg.UserHandler)

	// Private role routes
	pvtRole := private.Group("/roles")
	mapRolePrivateRoutes(pvtRole, cfg.RoleHandler)

	// Private application routes
	pvtApp := private.Group("/applications")
	mapAppPrivateRoutes(pvtApp, cfg.AppHandler)

	// Private permission routes
	pvtPerm := private.Group("/permissions")
	mapPermPrivateRoutes(pvtPerm, cfg.PermHandler)

	return nil
}

// mapAuthPublicRoutes maps public authentication routes
func mapAuthPublicRoutes(g *echo.Group, h *handler.AuthHandler) {
	g.POST("/login", h.Login())
}

// mapAuthPrivateRoutes maps private authentication routes
func mapAuthPrivateRoutes(g *echo.Group, h *handler.AuthHandler) {
	g.POST("/logout", h.Logout())
}

// mapUserPublicRoutes maps public user routes
func mapUserPublicRoutes(g *echo.Group, h *handler.UserHandler) {
	g.POST("", h.RegisterUser())
}

// mapUserPrivateRoutes maps private user routes
func mapUserPrivateRoutes(g *echo.Group, h *handler.UserHandler) {
	g.GET("/:id", h.GetUserProfile())
	g.GET("", h.GetUserList())
	g.PATCH("/:id", h.UpdateUserProfile())
	g.POST("/:id/applications/:app_id/roles", h.AssignUserRoles(), appMiddleware.RequirePermission("AUTHORIZER", "user.assign_roles"))
}

// mapRolePrivateRoutes maps private role routes
func mapRolePrivateRoutes(g *echo.Group, h *handler.RoleHandler) {
	g.POST("", h.Create(), appMiddleware.RequirePermission("AUTHORIZER", "role.create"))
	g.POST("/:id/permissions", h.GrantRolePermissions())
}

// mapAppPrivateRoutes maps private application routes
func mapAppPrivateRoutes(g *echo.Group, h *handler.AppHandler) {
	g.POST("", h.Create(), appMiddleware.RequirePermission("AUTHORIZER", "application.create"))
}

// mapPermPrivateRoutes maps private permission routes
func mapPermPrivateRoutes(g *echo.Group, h *handler.PermHandler) {
	g.POST("/sync", h.Sync(), appMiddleware.RequirePermission("AUTHORIZER", "permission.sync"))
}
