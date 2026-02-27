package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mafzaidi/authorizer/internal/delivery/http/middleware/pwd"
	"github.com/mafzaidi/authorizer/internal/domain/entity"
	"github.com/mafzaidi/authorizer/internal/domain/repository"
	"github.com/mafzaidi/authorizer/internal/domain/service"
	"github.com/mafzaidi/authorizer/pkg/idgen"
)

type userUsecase struct {
	repo         repository.UserRepository
	roleRepo     repository.RoleRepository
	userRoleRepo repository.UserRoleRepository
	logger       service.Logger
}

func NewUserUsecase(
	repo repository.UserRepository,
	roleRepo repository.RoleRepository,
	userRoleRepo repository.UserRoleRepository,
	logger service.Logger,
) Usecase {
	return &userUsecase{
		repo:         repo,
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
		logger:       logger,
	}
}

func (uc *userUsecase) Register(ctx context.Context, in *RegisterInput) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if len(in.Password) < 8 {
		uc.logger.Warn("Registration failed: password too short", service.Fields{
			"email": in.Email,
		})
		return errors.New("password must be at least 8 characters")
	}

	if in.Email == "" || in.Password == "" {
		uc.logger.Warn("Registration failed: missing required fields", service.Fields{})
		return errors.New("email and password required")
	}

	existingUser, _ := uc.repo.GetByEmail(ctx, in.Email)
	if existingUser != nil {
		uc.logger.Warn("Registration failed: email already exists", service.Fields{
			"email": in.Email,
		})
		return errors.New("email already exists")
	}

	hashedPassword, err := pwd.Hash(in.Password)
	if err != nil {
		uc.logger.Error("Failed to hash password", service.Fields{
			"email": in.Email,
			"error": err.Error(),
		})
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

	err = uc.repo.Create(ctx, user)
	if err != nil {
		uc.logger.Error("Failed to create user", service.Fields{
			"email": in.Email,
			"error": err.Error(),
		})
		return err
	}

	uc.logger.Info("User registered successfully", service.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	})

	return nil
}

func (uc *userUsecase) GetDetail(ctx context.Context, userID string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if userID == "" {
		uc.logger.Warn("GetDetail failed: userID is required", service.Fields{})
		return nil, errors.New("userID is required")
	}

	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get user detail", service.Fields{
			"user_id": userID,
			"error":   err.Error(),
		})
		return nil, fmt.Errorf("failed: %w", err)
	}

	return user, nil
}

func (uc *userUsecase) UpdateData(ctx context.Context, userID string, in *UpdateInput) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if userID == "" {
		uc.logger.Warn("UpdateData failed: userID is required", service.Fields{})
		return errors.New("userID is required")
	}

	existingUser, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get user for update", service.Fields{
			"user_id": userID,
			"error":   err.Error(),
		})
		return fmt.Errorf("failed: %w", err)
	}

	updatedUser := &entity.User{
		ID:       existingUser.ID,
		FullName: in.FullName,
		Phone:    &in.Phone,
	}

	err = uc.repo.Update(ctx, updatedUser)
	if err != nil {
		uc.logger.Error("Failed to update user", service.Fields{
			"user_id": userID,
			"error":   err.Error(),
		})
		return err
	}

	uc.logger.Info("User updated successfully", service.Fields{
		"user_id": userID,
	})

	return nil
}

func (uc *userUsecase) GetList(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if limit == 0 {
		limit = 50
	}

	users, err := uc.repo.List(ctx, limit, offset)
	if err != nil {
		uc.logger.Error("Failed to fetch users", service.Fields{
			"limit":  limit,
			"offset": offset,
			"error":  err.Error(),
		})
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	return users, nil
}

func (uc *userUsecase) AssignRoles(ctx context.Context, userID, appID string, roles []string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if userID == "" {
		uc.logger.Warn("AssignRoles failed: userID is required", service.Fields{})
		return errors.New("userID is required")
	}

	if appID == "" {
		uc.logger.Warn("AssignRoles failed: appID is required", service.Fields{
			"user_id": userID,
		})
		return errors.New("appID is required")
	}

	if len(roles) == 0 {
		uc.logger.Warn("AssignRoles failed: roles is required", service.Fields{
			"user_id": userID,
			"app_id":  appID,
		})
		return errors.New("roles is required")
	}

	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get user for role assignment", service.Fields{
			"user_id": userID,
			"error":   err.Error(),
		})
		return fmt.Errorf("failed: %w", err)
	}

	var roleIDs []string
	for _, v := range roles {
		role, err := uc.roleRepo.GetByAppAndCode(ctx, appID, v)
		if err != nil {
			uc.logger.Error("Failed to get role", service.Fields{
				"user_id":   userID,
				"app_id":    appID,
				"role_code": v,
				"error":     err.Error(),
			})
			return fmt.Errorf("failed: %w", err)
		}
		roleIDs = append(roleIDs, role.ID)
	}

	err = uc.userRoleRepo.Replace(ctx, user.ID, roleIDs)
	if err != nil {
		uc.logger.Error("Failed to assign roles", service.Fields{
			"user_id": userID,
			"app_id":  appID,
			"roles":   roles,
			"error":   err.Error(),
		})
		return err
	}

	uc.logger.Info("Roles assigned successfully", service.Fields{
		"user_id": userID,
		"app_id":  appID,
		"roles":   roles,
	})

	return nil
}
