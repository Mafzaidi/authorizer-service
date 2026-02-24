package v1

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
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
}

func NewPermHandler(uc perm.Usecase) *PermHandler {
	return &PermHandler{
		permUC: uc,
	}
}

func (h *PermHandler) Sync() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &SyncPermissionRequest{}

		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
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

		if err := h.permUC.SyncPermissions(c.Request().Context(), in); err != nil {
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		return response.SuccesHandler(c, &response.Response{
			Message: "permissions created successfully",
		})
	}
}
