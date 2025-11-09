package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	mid "github.com/labstack/echo/v4/middleware"
	httpHandler "localdev.me/authorizer/internal/delivery/http"
	"localdev.me/authorizer/internal/delivery/http/middleware"
	"localdev.me/authorizer/internal/domain/repository"
	"localdev.me/authorizer/internal/usecase"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	e.Use(mid.CORSWithConfig(mid.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
		AllowMethods:     []string{echo.POST, echo.GET, echo.OPTIONS},
	}))

	v1 := e.Group("/api/v1")
	userPublic := v1.Group("/user")
	health := v1.Group("/health")

	private := e.Group("/private/api/v1")
	private.Use(middleware.JWTAuthMiddleware)

	userRepo := repository.NewUserRepositoryPGX(s.db.Pool)
	userUC := usecase.NewUserUseCase(userRepo)
	userHandler := httpHandler.NewUserHandler(userUC)
	httpHandler.MapPublicRoutes(userPublic, userHandler, s.cfg)

	health.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct {
			Message string `json:"message"`
			Status  int    `json:"status"`
		}{
			Message: "Healthy",
			Status:  http.StatusOK,
		})
	})

	return nil
}
