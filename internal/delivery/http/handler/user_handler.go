package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mafzaidi/authorizer/internal/infrastructure/logger"
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
	logger *logger.Logger
}

func NewUserHandler(uc user.Usecase, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		userUC: uc,
		logger: logger,
	}
}

func (h *UserHandler) RegisterUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &RegisterUserRequest{}

		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			h.logger.Warn("Failed to decode register user request", logger.Fields{
				"error": err.Error(),
			})
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
			h.logger.Error("Failed to register user", logger.Fields{
				"email": req.Email,
				"error": err.Error(),
			})
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		h.logger.Info("User registered successfully", logger.Fields{
			"email":    req.Email,
			"username": req.Username,
		})

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
			h.logger.Warn("Failed to decode create user request", logger.Fields{
				"error": err.Error(),
			})
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
			h.logger.Error("Failed to create user", logger.Fields{
				"email": req.Email,
				"error": err.Error(),
			})
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		h.logger.Info("User created successfully", logger.Fields{
			"email":    req.Email,
			"username": req.Username,
		})

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
			h.logger.Warn("Failed to get user profile", logger.Fields{
				"user_id": userID,
				"error":   err.Error(),
			})
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

		h.logger.Info("User profile retrieved successfully", logger.Fields{
			"user_id": userID,
		})

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
			h.logger.Warn("Failed to decode update user request", logger.Fields{
				"user_id": userID,
				"error":   err.Error(),
			})
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		in := &user.UpdateInput{
			FullName: req.FullName,
			Phone:    req.Phone,
		}

		if err := h.userUC.UpdateData(c.Request().Context(), userID, in); err != nil {
			h.logger.Error("Failed to update user profile", logger.Fields{
				"user_id": userID,
				"error":   err.Error(),
			})
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		h.logger.Info("User profile updated successfully", logger.Fields{
			"user_id": userID,
		})

		return response.SuccesHandler(c, &response.Response{
			Message: "user data updated successfully",
		})
	}
}

func (h *UserHandler) GetUserList() echo.HandlerFunc {
	return func(c echo.Context) error {
		query := GetUserListQuery{}

		if err := c.Bind(&query); err != nil {
			h.logger.Warn("Failed to bind user list query", logger.Fields{
				"error": err.Error(),
			})
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
			h.logger.Error("Failed to get user list", logger.Fields{
				"page":   page,
				"limit":  limit,
				"offset": offset,
				"error":  err.Error(),
			})
			return response.ErrorHandler(c, http.StatusNotFound, "NotFound", err.Error())
		}

		h.logger.Info("User list retrieved successfully", logger.Fields{
			"page":  page,
			"limit": limit,
			"count": len(users),
		})

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
			h.logger.Warn("Failed to decode assign user roles request", logger.Fields{
				"user_id": userID,
				"app_id":  appID,
				"error":   err.Error(),
			})
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		roles := req.Roles

		if err := h.userUC.AssignRoles(c.Request().Context(), userID, appID, roles); err != nil {
			h.logger.Error("Failed to assign user roles", logger.Fields{
				"user_id": userID,
				"app_id":  appID,
				"roles":   roles,
				"error":   err.Error(),
			})
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		h.logger.Info("User roles assigned successfully", logger.Fields{
			"user_id": userID,
			"app_id":  appID,
			"roles":   roles,
		})

		return response.SuccesHandler(c, &response.Response{
			Message: "roles assigned successfully",
		})
	}
}
