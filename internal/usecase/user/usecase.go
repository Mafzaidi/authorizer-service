package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mafzaidi/authorizer/internal/delivery/http/middleware/pwd"
	"github.com/mafzaidi/authorizer/internal/domain/entity"
	"github.com/mafzaidi/authorizer/internal/domain/repository"
	"github.com/mafzaidi/authorizer/pkg/idgen"
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

func (uc *userUsecase) Register(ctx context.Context, in *RegisterInput) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if len(in.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if in.Email == "" || in.Password == "" {
		return errors.New("email and password required")
	}

	existingUser, _ := uc.repo.GetByEmail(ctx, in.Email)
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
		Phone:         &in.Phone,
		Password:      hashedPassword,
		Email:         in.Email,
		IsActive:      true,
		EmailVerified: false,
		PhoneVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return uc.repo.Create(ctx, user)
}

func (uc *userUsecase) GetDetail(ctx context.Context, userID string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if userID == "" {
		return nil, errors.New("userID is required")
	}

	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed: %w", err)
	}

	return user, nil
}

func (uc *userUsecase) UpdateData(ctx context.Context, userID string, in *UpdateInput) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if userID == "" {
		return errors.New("userID is required")
	}

	existingUser, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed: %w", err)
	}

	updatedUser := &entity.User{
		ID:       existingUser.ID,
		FullName: in.FullName,
		Phone:    &in.Phone,
	}

	return uc.repo.Update(ctx, updatedUser)
}

func (uc *userUsecase) GetList(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if limit == 0 {
		limit = 50
	}

	users, err := uc.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	return users, nil
}

func (uc *userUsecase) AssignRoles(ctx context.Context, userID, appID string, roles []string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if userID == "" {
		return errors.New("userID is required")
	}

	if appID == "" {
		return errors.New("appID is required")
	}

	if len(roles) == 0 {
		return errors.New("roles is required")
	}

	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed: %w", err)
	}

	var roleIDs []string
	for _, v := range roles {
		role, err := uc.roleRepo.GetByAppAndCode(ctx, appID, v)
		if err != nil {
			return fmt.Errorf("failed: %w", err)
		}
		roleIDs = append(roleIDs, role.ID)
	}

	return uc.userRoleRepo.Replace(ctx, user.ID, roleIDs)
}
