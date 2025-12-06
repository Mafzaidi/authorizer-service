package user

import (
	"errors"
	"fmt"
	"time"

	"localdev.me/authorizer/internal/delivery/http/middleware/pwd"
	"localdev.me/authorizer/internal/domain/entity"
	"localdev.me/authorizer/internal/domain/repository"
	"localdev.me/authorizer/pkg/idgen"
)

type userUsecase struct {
	repo         repository.UserRepository
	roleRepo     repository.RoleRepository
	userRoleRepo repository.UserRoleRepository
}

func NewUserUsecase(
	repo repository.UserRepository,
	roleRepo repository.RoleRepository,
	userRoleRepo repository.UserRoleRepository,
) Usecase {
	return &userUsecase{
		repo:         repo,
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
	}
}

func (u *userUsecase) Register(in *RegisterInput) error {
	if len(in.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if in.Email == "" || in.Password == "" {
		return errors.New("email and password required")
	}

	existingUser, _ := u.repo.GetByEmail(in.Email)
	if existingUser != nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := pwd.Hash(in.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user := &entity.User{
		ID:            idgen.NewUUIDv7(),
		Username:      in.Username,
		FullName:      in.FullName,
		Phone:         in.Phone,
		Password:      hashedPassword,
		Email:         in.Email,
		IsActive:      true,
		EmailVerified: false,
		PhoneVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return u.repo.Create(user)
}

func (u *userUsecase) GetDetail(userID string) (*entity.User, error) {
	if userID == "" {
		return nil, errors.New("userID is required")
	}

	user, err := u.repo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed: %w", err)
	}

	return user, nil
}

func (u *userUsecase) UpdateData(userID string, in *UpdateInput) error {
	if userID == "" {
		return errors.New("userID is required")
	}

	existingUser, err := u.repo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed: %w", err)
	}

	updatedUser := &entity.User{
		ID:       existingUser.ID,
		FullName: in.FullName,
		Phone:    in.Phone,
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
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	return users, nil
}

func (u *userUsecase) AssignRoles(userID, appID string, roles []string) error {

	if userID == "" {
		return errors.New("userID is required")
	}

	if appID == "" {
		return errors.New("appID is required")
	}

	if len(roles) == 0 {
		return errors.New("roles is required")
	}

	user, err := u.repo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed: %w", err)
	}

	var roleIDs []string
	for _, v := range roles {
		role, err := u.roleRepo.GetByAppAndCode(appID, v)
		if err != nil {
			return fmt.Errorf("failed: %w", err)
		}
		roleIDs = append(roleIDs, role.ID)
	}

	return u.userRoleRepo.Replace(user.ID, roleIDs)
}
