package http

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"localdev.me/authorizer/internal/usecase"
	"localdev.me/authorizer/pkg/response"
)

type User interface {
	Register() echo.HandlerFunc
}

type UserHandler struct {
	userUC usecase.UserUseCase
}

func NewUserHandler(uc usecase.UserUseCase) User {
	return &UserHandler{
		userUC: uc,
	}
}
func (h *UserHandler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		pl := &usecase.RegisterRequest{}

		if err := json.NewDecoder(c.Request().Body).Decode(&pl); err != nil {
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		if err := h.userUC.Register(pl); err != nil {
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		return response.SuccesHandler(c, &response.Response{
			Message: "user registered successfully",
		})
	}
}
