package user

import "localdev.me/authorizer/internal/domain/entity"

type Usecase interface {
	Register(req *RegisterInput) error
	GetDetail(userID string) (*entity.User, error)
	UpdateData(userID string, input *UpdateInput) error
	GetList(limit, offset int) ([]*entity.User, error)
	AssignRoles(userID, appID string, roles []string) error
}
