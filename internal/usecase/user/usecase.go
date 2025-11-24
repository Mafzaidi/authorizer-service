package user

import (
	"errors"
	"time"

	"localdev.me/authorizer/internal/delivery/http/middleware/pwd"
	"localdev.me/authorizer/internal/domain/entity"
	"localdev.me/authorizer/internal/domain/repository"
	"localdev.me/authorizer/pkg/idgen"
)

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) Usecase {
	return &userUsecase{
		repo: repo,
	}
}

func (u *userUsecase) Register(input *RegisterInput) error {
	if len(input.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if input.Email == "" || input.Password == "" {
		return errors.New("email and password required")
	}

	existing, _ := u.repo.GetByEmail(input.Email)
	if existing != nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := pwd.Hash(input.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	newUser := &entity.User{
		ID:            idgen.NewUUIDv7(),
		Username:      input.Username,
		FullName:      input.FullName,
		Phone:         input.Phone,
		Password:      hashedPassword,
		Email:         input.Email,
		IsActive:      true,
		EmailVerified: false,
		PhoneVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return u.repo.Create(newUser)
}

func (u *userUsecase) GetDetail(id string) (*entity.User, error) {
	if id == "" {
		return nil, errors.New("userID is required")
	}

	user, err := u.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (u *userUsecase) UpdateData(id string, input *UpddateInput) error {
	if id == "" {
		return errors.New("userID is required")
	}

	user, err := u.repo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	updatedUser := &entity.User{
		ID:       user.ID,
		FullName: input.FullName,
		Phone:    input.Phone,
	}

	return u.repo.Update(updatedUser)
}

func (u *userUsecase) GetList(limit, offset int) ([]*entity.User, error) {

	if limit == 0 {
		limit = 50
	}

	if offset == 0 {
		offset = 10
	}

	users, err := u.repo.List(limit, offset)
	if err != nil {
		return nil, errors.New("no users found")
	}

	return users, nil
}
