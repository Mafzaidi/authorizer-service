package handler

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mafzaidi/authorizer/internal/domain/service"
	perm "github.com/mafzaidi/authorizer/internal/usecase/permission"
	"github.com/mafzaidi/authorizer/pkg/response"
)

type PermissionRequest struct {
	Code        string `json:"code" validate:"required"`
	Description string `json:"description"`
}

type SyncPermissionRequest struct {
	Application string `json:"application" validate:"required"`
	Permissions []*PermissionRequest
	Version     int `json:"version" validate:"required"`
}

type PermHandler struct {
	permUC perm.Usecase
	logger service.Logger
}

func NewPermHandler(uc perm.Usecase, logger service.Logger) *PermHandler {
	return &PermHandler{
		permUC: uc,
		logger: logger,
	}
}

func (h *PermHandler) Sync() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &SyncPermissionRequest{}

		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			h.logger.Warn("Failed to decode sync permission request", service.Fields{
				"error": err.Error(),
			})
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		var perms []*perm.PermissionsInput
		for _, v := range req.Permissions {
			perm := &perm.PermissionsInput{
				Code:        v.Code,
				Description: v.Description,
			}
			perms = append(perms, perm)
		}

		in := &perm.SyncInput{
			AppCode:     req.Application,
			Permissions: perms,
			Version:     req.Version,
		}

		h.logger.Info("Syncing permissions", service.Fields{
			"app_code":          req.Application,
			"permissions_count": len(perms),
			"version":           req.Version,
		})

		if err := h.permUC.SyncPermissions(c.Request().Context(), in); err != nil {
			h.logger.Error("Failed to sync permissions", service.Fields{
				"app_code": req.Application,
				"error":    err.Error(),
			})
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		h.logger.Info("Permissions synced successfully", service.Fields{
			"app_code":          req.Application,
			"permissions_count": len(perms),
		})

		return response.SuccesHandler(c, &response.Response{
			Message: "permissions created successfully",
		})
	}
}
