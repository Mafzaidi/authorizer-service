package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/mafzaidi/authorizer/internal/infrastructure/logger"
	"github.com/mafzaidi/authorizer/pkg/response"
)

type HealthHandler struct {
	logger *logger.Logger
}

func NewHealthHandler(logger *logger.Logger) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

func (h *HealthHandler) Check() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.logger.Debug("Health check endpoint called", logger.Fields{})

		return response.SuccesHandler(c, &response.Response{
			Message: "OK",
		})
	}
}
