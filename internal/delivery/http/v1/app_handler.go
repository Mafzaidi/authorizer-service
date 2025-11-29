package v1

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	app "localdev.me/authorizer/internal/usecase/application"
	"localdev.me/authorizer/pkg/response"
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
		pl := &CreateAppRequest{}

		if err := json.NewDecoder(c.Request().Body).Decode(&pl); err != nil {
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		input := &app.CreateInput{
			Code:        pl.Code,
			Name:        pl.Name,
			Description: pl.Description,
		}

		if err := h.appUC.Create(input); err != nil {
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		return response.SuccesHandler(c, &response.Response{
			Message: "application created successfully",
		})
	}
}
