package v1

import (
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mafzaidi/authorizer/config"
	"github.com/mafzaidi/authorizer/internal/usecase/auth"
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
		Keys []JWK `json:"keys"`
	}

	JWK struct {
		Kty string `json:"kty"`
		Use string `json:"use"`
		Alg string `json:"alg"`
		Kid string `json:"kid"`
		N   string `json:"n"`
		E   string `json:"e"`
	}

	Authorization struct {
		App         string   `json:"app"`
		Roles       []string `json:"roles"`
		Permissions []string `json:"permissions"`
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

		data, err := h.authUC.Login(c.Request().Context(), req.Application, req.Email, req.Password, validToken, h.cfg)
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

func (h *AuthHandler) GetJWKS() echo.HandlerFunc {
	return func(c echo.Context) error {
		pub := h.cfg.JWT.PublicKey

		n := base64.RawURLEncoding.EncodeToString(pub.N.Bytes())
		e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes())

		resp := JWKSResponse{
			Keys: []JWK{
				{
					Kty: "RSA",
					Use: "sig",
					Alg: "RS256",
					Kid: h.cfg.JWT.KeyID,
					N:   n,
					E:   e,
				},
			},
		}

		return c.JSON(http.StatusOK, resp)
	}
}
