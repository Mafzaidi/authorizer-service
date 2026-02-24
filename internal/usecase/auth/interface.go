package auth

import (
	"context"

	"github.com/mafzaidi/authorizer/config"
)

type Usecase interface {
	Login(ctx context.Context, application, email, password, validToken string, conf *config.Config) (*UserToken, error)
}
