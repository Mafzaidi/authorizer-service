package v1

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mafzaidi/authorizer/internal/usecase/role"
	"github.com/mafzaidi/authorizer/pkg/response"
)

type (
	CreateRoleRequest struct {
		AppID       string `json:"application_id"`
		Code        string `json:"code" validate:"required"`
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"required"`
	}

	GrantRolePermissionsRequest struct {
		Perms []string `json:"permissions" validate:"required"`
	}
)

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
		req := &CreateRoleRequest{}

		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		in := &role.CreateInput{
			AppID:       req.AppID,
			Code:        req.Code,
			Name:        req.Name,
			Description: req.Description,
		}

		if err := h.roleUC.Create(c.Request().Context(), in); err != nil {
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		return response.SuccesHandler(c, &response.Response{
			Message: "role created successfully",
		})
	}
}

func (h *RoleHandler) GrantRolePermissions() echo.HandlerFunc {
	return func(c echo.Context) error {
		roleID := c.Param("id")
		req := &GrantRolePermissionsRequest{}
		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		perms := req.Perms

		if err := h.roleUC.GrantPerms(c.Request().Context(), roleID, perms); err != nil {
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		return response.SuccesHandler(c, &response.Response{
			Message: "permissions granted successfully",
		})
	}
}
