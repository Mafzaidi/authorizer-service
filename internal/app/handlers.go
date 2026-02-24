package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
	mid "github.com/labstack/echo/v4/middleware"
	"github.com/mafzaidi/authorizer/internal/app/resource"
	"github.com/mafzaidi/authorizer/internal/delivery/http/middleware"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	e.Use(mid.CORSWithConfig(mid.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
		AllowMethods:     []string{echo.POST, echo.GET, echo.OPTIONS},
	}))

	auth := resource.NewAuth(s.postgreDB.Pool, s.cfg, s.redis.Client)
	user := resource.NewUser(s.postgreDB.Pool)
	role := resource.NewRole(s.postgreDB.Pool)
	app := resource.NewApp(s.postgreDB.Pool)
	perm := resource.NewPerm(s.postgreDB.Pool)

	v1 := e.Group("authorizer/v1")
	public := v1.Group("")

	pblAuth := public.Group("/auth")
	pblUser := public.Group("/users")
	pblHealth := public.Group("/health")

	MapAuthPublicRoutes(pblAuth, auth.Handler)
	MapUserPublicRoutes(pblUser, user.Handler, s.cfg)

	private := v1.Group("")
	private.Use(middleware.JWTAuthMiddleware)
	pvtAuth := private.Group("/auth")
	pvtUser := private.Group("/users")
	pvtRole := private.Group("/roles")
	pvtApp := private.Group("/applications")
	pvtPerm := private.Group("/permissions")

	MapAuthPrivateRoutes(pvtAuth, auth.Handler)
	MapUserPrivateRoutes(pvtUser, user.Handler)
	MapRolePrivateRoutes(pvtRole, role.Handler)
	MapAppPrivateRoutes(pvtApp, app.Handler)
	MapPermPrivateRoutes(pvtPerm, perm.Handler)

	pblHealth.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"message": "Healthy",
			"status":  http.StatusOK,
		})
	})

	e.GET("/.well-known/jwks.json", auth.Handler.GetJWKS())

	return nil
}
