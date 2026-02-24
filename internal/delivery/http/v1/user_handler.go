package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mafzaidi/authorizer/internal/usecase/user"
	"github.com/mafzaidi/authorizer/pkg/response"
)

type (
	RegisterUserRequest struct {
		Username string `json:"username" validate:"required"`
		FullName string `json:"full_name" validate:"required"`
		Phone    string `json:"phone" validate:"required"`
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	UpdateUserRequest struct {
		FullName string `json:"full_name" validate:"required"`
		Phone    string `json:"phone" validate:"required"`
	}

	GetUserProfileResponse struct {
		ID          string    `json:"id"`
		Username    string    `json:"username"`
		Fullname    string    `json:"full_name"`
		PhoneNumber *string   `json:"phone_number"`
		Email       string    `json:"email"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	GetUserListQuery struct {
		Page  int `query:"page"`
		Limit int `query:"limit"`
	}

	AssignUserRoleRequest struct {
		Roles []string `json:"roles" validate:"required"`
	}
)

type UserHandler struct {
	userUC user.Usecase
}

func NewUserHandler(uc user.Usecase) *UserHandler {
	return &UserHandler{
		userUC: uc,
	}
}

func (h *UserHandler) RegisterUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &RegisterUserRequest{}

		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		in := &user.RegisterInput{
			Username: req.Username,
			FullName: req.FullName,
			Phone:    req.Phone,
			Email:    req.Email,
			Password: req.Password,
		}

		if err := h.userUC.Register(c.Request().Context(), in); err != nil {
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		return response.SuccesHandler(c, &response.Response{
			Message: "user registered successfully",
		})
	}
}

func (h *UserHandler) CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		// claims := middleware.GetUserFromContext(c)

		// if !middleware.HasRole("superadmin", claims.Authorization) {
		// 	return response.ErrorHandler(c, http.StatusForbidden, "Forbidden", "you don't have access to this route")
		// }

		req := &RegisterUserRequest{}

		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		in := &user.RegisterInput{
			Username: req.Username,
			FullName: req.FullName,
			Phone:    req.Phone,
			Email:    req.Email,
			Password: req.Password,
		}

		if err := h.userUC.Register(c.Request().Context(), in); err != nil {
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		return response.SuccesHandler(c, &response.Response{
			Message: "user registered successfully",
		})
	}
}

func (h *UserHandler) GetUserProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Param("id")
		// claims := middleware.GetUserFromContext(c)

		// if !middleware.HasRole("superadmin", claims.Roles) && claims.UserID != userID {
		// 	return response.ErrorHandler(c, http.StatusForbidden, "Forbidden", "you don't have access to this route")
		// }

		user, err := h.userUC.GetDetail(c.Request().Context(), userID)
		if err != nil {
			return response.ErrorHandler(c, http.StatusNotFound, "NotFound", err.Error())
		}

		resp := &GetUserProfileResponse{
			ID:          user.ID,
			Username:    user.Username,
			Fullname:    user.FullName,
			PhoneNumber: user.Phone,
			Email:       user.Email,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		}

		return response.SuccesHandler(c, &response.Response{
			Message: "get user data successfully",
			Data:    resp,
		})

	}
}

func (h *UserHandler) UpdateUserProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Param("id")
		req := &UpdateUserRequest{}
		// claims := middleware.GetUserFromContext(c)

		// if !middleware.HasRole("superadmin", claims.Roles) && claims.UserID != userID {
		// 	return response.ErrorHandler(c, http.StatusForbidden, "Forbidden", "you don't have access to this route")
		// }

		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		in := &user.UpdateInput{
			FullName: req.FullName,
			Phone:    req.Phone,
		}

		if err := h.userUC.UpdateData(c.Request().Context(), userID, in); err != nil {
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		return response.SuccesHandler(c, &response.Response{
			Message: "user data updated successfully",
		})
	}
}

func (h *UserHandler) GetUserList() echo.HandlerFunc {
	return func(c echo.Context) error {
		query := GetUserListQuery{}

		if err := c.Bind(&query); err != nil {
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		// claims := middleware.GetUserFromContext(c)

		// if !middleware.HasRole("superadmin", claims.Roles) {
		// 	return response.ErrorHandler(c, http.StatusForbidden, "Forbidden", "you don't have access to this route")
		// }

		page := query.Page
		if page <= 0 {
			page = 1
		}
		limit := query.Limit
		if limit <= 0 {
			limit = 50 // default limit
		}
		offset := (page - 1) * limit

		users, err := h.userUC.GetList(c.Request().Context(), limit, offset)
		if err != nil {
			return response.ErrorHandler(c, http.StatusNotFound, "NotFound", err.Error())
		}

		return response.SuccesHandler(c, &response.Response{
			Message: "OK",
			Data:    users,
		})

	}
}

func (h *UserHandler) AssignUserRoles() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Param("id")
		appID := c.Param("app_id")
		req := &AssignUserRoleRequest{}
		// claims := middleware.GetUserFromContext(c)

		// if claims.UserID != userID {
		// 	return response.ErrorHandler(c, http.StatusForbidden, "Forbidden", "you don't have access to this route")
		// }

		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		roles := req.Roles

		if err := h.userUC.AssignRoles(c.Request().Context(), userID, appID, roles); err != nil {
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		return response.SuccesHandler(c, &response.Response{
			Message: "roles assigned successfully",
		})
	}
}
