package handler

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mafzaidi/authorizer/internal/domain/service"
	app "github.com/mafzaidi/authorizer/internal/usecase/application"
	"github.com/mafzaidi/authorizer/pkg/response"
)

type CreateAppRequest struct {
	Code        string `json:"code" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type AppHandler struct {
	appUC  app.Usecase
	logger service.Logger
}

func NewAppHandler(uc app.Usecase, logger service.Logger) *AppHandler {
	return &AppHandler{
		appUC:  uc,
		logger: logger,
	}
}

func (h *AppHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &CreateAppRequest{}

		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			h.logger.Warn("Failed to decode application creation request", service.Fields{
				"error": err.Error(),
			})
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		in := &app.CreateInput{
			Code:        req.Code,
			Name:        req.Name,
			Description: req.Description,
		}

		if err := h.appUC.Create(c.Request().Context(), in); err != nil {
			h.logger.Error("Failed to create application", service.Fields{
				"code":  req.Code,
				"error": err.Error(),
			})
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		h.logger.Info("Application created successfully via HTTP", service.Fields{
			"code": req.Code,
			"name": req.Name,
		})

		return response.SuccesHandler(c, &response.Response{
			Message: "application created successfully",
		})
	}
}
