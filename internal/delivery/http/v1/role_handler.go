package v1

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"localdev.me/authorizer/internal/usecase/role"
	"localdev.me/authorizer/pkg/response"
)

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Application string `json:"application" validate:"required"`
}

type RoleHandler struct {
	roleUC role.Usecase
}

func NewRoleHandler(uc role.Usecase) *RoleHandler {
	return &RoleHandler{
		roleUC: uc,
	}
}

func (h *RoleHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		pl := &CreateRoleRequest{}

		if err := json.NewDecoder(c.Request().Body).Decode(&pl); err != nil {
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		name := pl.Name
		description := pl.Description
		application := pl.Application

		if err := h.roleUC.Create(name, description, application); err != nil {
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		return response.SuccesHandler(c, &response.Response{
			Message: "role created successfully",
		})
	}
}
