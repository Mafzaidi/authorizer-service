package v1

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	app "github.com/mafzaidi/authorizer/internal/usecase/application"
	"github.com/mafzaidi/authorizer/pkg/response"
)

type CreateAppRequest struct {
	Code        string `json:"code" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type AppHandler struct {
	appUC app.Usecase
}

func NewAppHandler(uc app.Usecase) *AppHandler {
	return &AppHandler{
		appUC: uc,
	}
}

func (h *AppHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &CreateAppRequest{}

		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		in := &app.CreateInput{
			Code:        req.Code,
			Name:        req.Name,
			Description: req.Description,
		}

		if err := h.appUC.Create(c.Request().Context(), in); err != nil {
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		return response.SuccesHandler(c, &response.Response{
			Message: "application created successfully",
		})
	}
}
