package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mafzaidi/authorizer/internal/infrastructure/auth"
	"github.com/mafzaidi/authorizer/internal/infrastructure/config"
	"github.com/mafzaidi/authorizer/internal/infrastructure/logger"
	authUsecase "github.com/mafzaidi/authorizer/internal/usecase/auth"
	"github.com/mafzaidi/authorizer/pkg/response"
)

type (
	LoginRequest struct {
		Application string `json:"application"`
		Email       string `json:"email"`
		Password    string `json:"password"`
	}

	LoginResponse struct {
		Username    string `json:"username"`
		Fullname    string `json:"full_name"`
		AccessToken struct {
			Type      string    `json:"type"`
			Token     string    `json:"token"`
			ExpiresAt time.Time `json:"expires_at"`
		} `json:"access_token"`
		RefreshToken string `json:"refresh_token"`

		Authorization []Authorization `json:"authorization"`
	}

	JWKSResponse struct {
		Keys []auth.JWK `json:"keys"`
	}

	Authorization struct {
		App         string   `json:"app"`
		Roles       []string `json:"roles"`
		Permissions []string `json:"permissions"`
	}
)

type AuthHandler struct {
	authUC      authUsecase.Usecase
	jwksService auth.JWKSService
	cfg         *config.Config
	logger      *logger.Logger
}

func NewAuthHandler(
	authUC authUsecase.Usecase,
	jwksService auth.JWKSService,
	cfg *config.Config,
	logger *logger.Logger,
) *AuthHandler {
	return &AuthHandler{
		authUC:      authUC,
		jwksService: jwksService,
		cfg:         cfg,
		logger:      logger,
	}
}

func (h *AuthHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &LoginRequest{}
		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			h.logger.Warn("Failed to decode login request", logger.Fields{
				"error": err.Error(),
			})
			return response.ErrorHandler(c, http.StatusBadRequest, "BadRequest", err.Error())
		}

		var validToken string
		if cookie, err := c.Cookie("jwt_user_token"); err == nil {
			validToken = cookie.Value
		}

		data, err := h.authUC.Login(c.Request().Context(), req.Application, req.Email, req.Password, validToken, h.cfg)
		if err != nil {
			h.logger.Warn("Login failed", logger.Fields{
				"email": req.Email,
				"error": err.Error(),
			})
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

		var authorizations []Authorization

		for _, a := range data.Claims.Authorization {
			authorizations = append(authorizations, Authorization{
				App:         a.App,
				Roles:       a.Roles,
				Permissions: a.Permissions,
			})
		}

		resp := &LoginResponse{
			Username: data.Claims.Username,
			Fullname: data.User.FullName,
			AccessToken: struct {
				Type      string    `json:"type"`
				Token     string    `json:"token"`
				ExpiresAt time.Time `json:"expires_at"`
			}{
				Type:      "Bearer",
				Token:     data.Token,
				ExpiresAt: data.Claims.ExpiresAt.Time,
			},
			RefreshToken:  data.RefreshToken,
			Authorization: authorizations,
		}

		h.logger.Info("User logged in successfully", logger.Fields{
			"email":    req.Email,
			"username": data.Claims.Username,
		})

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

		h.logger.Info("User logged out successfully", logger.Fields{})

		return response.SuccesHandler(c, &response.Response{
			Message: "user logged out successfully",
		})
	}
}

func (h *AuthHandler) GetJWKS() echo.HandlerFunc {
	return func(c echo.Context) error {
		jwksResp, err := h.jwksService.GetJWKS(h.cfg.JWT.PublicKey, h.cfg.JWT.KeyID)
		if err != nil {
			h.logger.Error("Failed to generate JWKS", logger.Fields{
				"error": err.Error(),
			})
			return response.ErrorHandler(c, http.StatusInternalServerError, "InternalServerError", "failed to generate JWKS")
		}

		resp := JWKSResponse{
			Keys: jwksResp.Keys,
		}

		return c.JSON(http.StatusOK, resp)
	}
}
