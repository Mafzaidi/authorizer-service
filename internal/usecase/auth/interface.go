package auth

import "localdev.me/authorizer/config"

type Usecase interface {
	Login(application, email, password, validToken string, conf *config.Config) (*UserToken, error)
}
