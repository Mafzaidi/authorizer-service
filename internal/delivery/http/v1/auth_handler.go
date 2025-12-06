package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"localdev.me/authorizer/config"
	"localdev.me/authorizer/internal/usecase/auth"
	"localdev.me/authorizer/pkg/response"
)

type (
	LoginRequest struct {
		Application string `json:"application"`
		Email       string `json:"email"`
		Password    string `json:"password"`
	}

	LoginResponse struct {
		Username    string   `json:"username"`
		Fullname    string   `json:"full_name"`
		Roles       []string `json:"roles"`
		Permissions []string `json:"permission"`
		AccessToken struct {
			Type      string    `json:"type"`
			Token     string    `json:"token"`
			ExpiresAt time.Time `json:"expires_at"`
		} `json:"access_token"`
	}
)

type AuthHandler struct {
	authUC auth.Usecase
	cfg    *config.Config
}

func NewAuthHandler(uc auth.Usecase, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authUC: uc,
		cfg:    cfg,
	}
}

func (h *AuthHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &LoginRequest{}
		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		var validToken string
		if cookie, err := c.Cookie("jwt_user_token"); err == nil {
			validToken = cookie.Value
		}

		data, err := h.authUC.Login(req.Application, req.Email, req.Password, validToken, h.cfg)
		if err != nil {
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		newCookie := new(http.Cookie)
		newCookie.Name = "jwt_user_token"
		newCookie.Value = data.Token
		// newCookie.HttpOnly = true
		newCookie.Secure = true
		newCookie.SameSite = http.SameSiteNoneMode
		newCookie.Expires = data.Claims.ExpiresAt.Time
		newCookie.Path = "/"

		c.SetCookie(newCookie)

		resp := &LoginResponse{
			Username:    data.Claims.Username,
			Fullname:    data.User.FullName,
			Roles:       data.Claims.Roles,
			Permissions: data.Claims.Permissions,
			AccessToken: struct {
				Type      string    `json:"type"`
				Token     string    `json:"token"`
				ExpiresAt time.Time `json:"expires_at"`
			}{
				Type:      "Bearer",
				Token:     data.Token,
				ExpiresAt: data.Claims.ExpiresAt.Time,
			},
		}

		return response.SuccesHandler(c, &response.Response{
			Message: "user login successfully",
			Data:    resp,
		})
	}
}

func (h *AuthHandler) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		expiredCookie := &http.Cookie{
			Name:     "jwt_user_token",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		}
		c.SetCookie(expiredCookie)

		return response.SuccesHandler(c, &response.Response{
			Message: "user logged out successfully",
		})
	}
}
