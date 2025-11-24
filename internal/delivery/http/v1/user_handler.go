package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"localdev.me/authorizer/internal/delivery/http/middleware"
	"localdev.me/authorizer/internal/usecase/user"
	"localdev.me/authorizer/pkg/response"
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
		PhoneNumber string    `json:"phone_number"`
		Email       string    `json:"email"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	GetUserListQuery struct {
		Offset int `query:"page"`
		Limit  int `query:"limit"`
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
		pl := &RegisterUserRequest{}

		if err := json.NewDecoder(c.Request().Body).Decode(&pl); err != nil {
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		input := &user.RegisterInput{
			Username: pl.Username,
			FullName: pl.FullName,
			Phone:    pl.Phone,
			Email:    pl.Email,
			Password: pl.Password,
		}

		if err := h.userUC.Register(input); err != nil {
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		return response.SuccesHandler(c, &response.Response{
			Message: "user registered successfully",
		})
	}
}

func (h *UserHandler) CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := middleware.GetUserFromContext(c)

		if !middleware.HasRole("superadmin", claims.Roles) {
			return response.ErrorHandler(c, http.StatusForbidden, "Forbidden", "you don't have access to this route")
		}

		pl := &RegisterUserRequest{}

		if err := json.NewDecoder(c.Request().Body).Decode(&pl); err != nil {
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		input := &user.RegisterInput{
			Username: pl.Username,
			FullName: pl.FullName,
			Phone:    pl.Phone,
			Email:    pl.Email,
			Password: pl.Password,
		}

		if err := h.userUC.Register(input); err != nil {
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", err.Error())
		}

		return response.SuccesHandler(c, &response.Response{
			Message: "user registered successfully",
		})
	}
}

func (h *UserHandler) GetUserProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		ID := c.Param("id")
		claims := middleware.GetUserFromContext(c)

		if !middleware.HasRole("superadmin", claims.Roles) && claims.UserID != ID {
			return response.ErrorHandler(c, http.StatusForbidden, "Forbidden", "you don't have access to this route")
		}

		user, err := h.userUC.GetDetail(ID)
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
		ID := c.Param("id")
		pl := &UpdateUserRequest{}

		claims := middleware.GetUserFromContext(c)

		if !middleware.HasRole("superadmin", claims.Roles) && claims.UserID != ID {
			return response.ErrorHandler(c, http.StatusForbidden, "Forbidden", "you don't have access to this route")
		}

		input := &user.UpddateInput{
			FullName: pl.FullName,
			Phone:    pl.Phone,
		}

		if err := h.userUC.UpdateData(ID, input); err != nil {
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

		claims := middleware.GetUserFromContext(c)

		if !middleware.HasRole("superadmin", claims.Roles) {
			return response.ErrorHandler(c, http.StatusForbidden, "Forbidden", "you don't have access to this route")
		}

		users, err := h.userUC.GetList(query.Offset, query.Limit)
		if err != nil {
			return response.ErrorHandler(c, http.StatusNotFound, "NotFound", err.Error())
		}

		return response.SuccesHandler(c, &response.Response{
			Message: "OK",
			Data:    users,
		})

	}
}
