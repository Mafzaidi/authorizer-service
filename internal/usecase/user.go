package usecase

import (
	"errors"
	"time"

	"localdev.me/authorizer/internal/delivery/http/middleware/pwd"
	"localdev.me/authorizer/internal/domain/entity"
	"localdev.me/authorizer/internal/domain/repository"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserUseCase interface {
	Register(req *RegisterRequest) error
}

type userUC struct {
	repo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &userUC{
		repo: repo,
	}
}

func (u *userUC) Register(req *RegisterRequest) error {
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if req.Email == "" || req.Password == "" {
		return errors.New("email dan password wajib diisi")
	}

	existing, _ := u.repo.GetByEmail(req.Email)
	if existing == nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := pwd.Hash(req.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	newUser := &entity.User{
		ID:            uuid.NewString(),
		Username:      req.Username,
		FullName:      req.FullName,
		Phone:         req.Phone,
		PasswordHash:  hashedPassword,
		Email:         req.Email,
		Role:          "user",
		IsActive:      true,
		EmailVerified: false,
		PhoneVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return u.repo.Create(newUser)
}
