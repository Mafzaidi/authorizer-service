package user

import "localdev.me/authorizer/internal/domain/entity"

type Usecase interface {
	Register(req *RegisterInput) error
	GetDetail(id string) (*entity.User, error)
	UpdateData(id string, input *UpdateInput) error
	GetList(limit, offset int) ([]*entity.User, error)
}
