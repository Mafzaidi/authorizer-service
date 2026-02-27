package handler

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mafzaidi/authorizer/internal/infrastructure/logger"
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
	logger *logger.Logger
}

func NewRoleHandler(uc role.Usecase, logger *logger.Logger) *RoleHandler {
	return &RoleHandler{
		roleUC: uc,
		logger: logger,
	}
}

func (h *RoleHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &CreateRoleRequest{}

		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			h.logger.Warn("Failed to decode create role request", logger.Fields{
				"error": err.Error(),
			})
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		in := &role.CreateInput{
			AppID:       req.AppID,
			Code:        req.Code,
			Name:        req.Name,
			Description: req.Description,
		}

		if err := h.roleUC.Create(c.Request().Context(), in); err != nil {
			h.logger.Error("Failed to create role", logger.Fields{
				"app_id": req.AppID,
				"code":   req.Code,
				"error":  err.Error(),
			})
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		h.logger.Info("Role created successfully", logger.Fields{
			"app_id": req.AppID,
			"code":   req.Code,
			"name":   req.Name,
		})

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
			h.logger.Warn("Failed to decode grant role permissions request", logger.Fields{
				"role_id": roleID,
				"error":   err.Error(),
			})
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		perms := req.Perms

		if err := h.roleUC.GrantPerms(c.Request().Context(), roleID, perms); err != nil {
			h.logger.Error("Failed to grant role permissions", logger.Fields{
				"role_id":     roleID,
				"permissions": perms,
				"error":       err.Error(),
			})
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		h.logger.Info("Role permissions granted successfully", logger.Fields{
			"role_id":     roleID,
			"permissions": perms,
		})

		return response.SuccesHandler(c, &response.Response{
			Message: "permissions granted successfully",
		})
	}
}
