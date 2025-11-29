package v1

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	perm "localdev.me/authorizer/internal/usecase/permission"
	"localdev.me/authorizer/pkg/response"
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
		pl := &SyncPermissionRequest{}

		if err := json.NewDecoder(c.Request().Body).Decode(&pl); err != nil {
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		var perms []*perm.PermissionsInput
		for _, p := range pl.Permissions {
			perm := &perm.PermissionsInput{
				Code:        p.Code,
				Description: p.Description,
			}
			perms = append(perms, perm)
		}

		input := &perm.SyncInput{
			AppCode:     pl.Application,
			Permissions: perms,
			Version:     pl.Version,
		}

		if err := h.permUC.SyncPermissions(input); err != nil {
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		return response.SuccesHandler(c, &response.Response{
			Message: "role created successfully",
		})
	}
}
