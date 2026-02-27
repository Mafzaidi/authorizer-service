package application

import (
	"context"
	"errors"
	"time"

	"github.com/mafzaidi/authorizer/internal/domain/entity"
	"github.com/mafzaidi/authorizer/internal/domain/repository"
	"github.com/mafzaidi/authorizer/internal/domain/service"
	"github.com/mafzaidi/authorizer/pkg/idgen"
)

type appUsecase struct {
	repo   repository.AppRepository
	logger service.Logger
}

func NewAppUsecase(repo repository.AppRepository, logger service.Logger) Usecase {
	return &appUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *appUsecase) Create(ctx context.Context, in *CreateInput) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if in.Code == "" || in.Name == "" {
		uc.logger.Warn("Application creation failed: code and name required", service.Fields{})
		return errors.New("code and name is required")
	}

	existingApp, _ := uc.repo.GetByCode(ctx, in.Code)
	if existingApp != nil {
		uc.logger.Warn("Application creation failed: application already exists", service.Fields{
			"code": in.Code,
		})
		return errors.New("application already exists")
	}

	app := &entity.Application{
		ID:          idgen.NewUUIDv7(),
		Code:        in.Code,
		Description: in.Description,
		Name:        in.Name,
		Metadata:    in.Metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := uc.repo.Create(ctx, app)
	if err != nil {
		uc.logger.Error("Failed to create application", service.Fields{
			"code":  in.Code,
			"error": err.Error(),
		})
		return err
	}

	uc.logger.Info("Application created successfully", service.Fields{
		"id":   app.ID,
		"code": app.Code,
		"name": app.Name,
	})

	return nil
}
