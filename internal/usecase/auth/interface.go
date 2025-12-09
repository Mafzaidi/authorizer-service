package auth

import (
	"context"

	"localdev.me/authorizer/config"
)

type Usecase interface {
	Login(ctx context.Context, application, email, password, validToken string, conf *config.Config) (*UserToken, error)
}
