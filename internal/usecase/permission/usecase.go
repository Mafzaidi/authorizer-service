package permission

import (
	"errors"
	"fmt"
	"time"

	"localdev.me/authorizer/internal/domain/entity"
	"localdev.me/authorizer/internal/domain/repository"
	"localdev.me/authorizer/pkg/idgen"
)

type permUsecase struct {
	permRepo repository.PermRepository
	appRepo  repository.AppRepository
}

func NewPermUsecase(
	permRepo repository.PermRepository,
	appRepo repository.AppRepository,
) Usecase {
	return &permUsecase{
		permRepo: permRepo,
		appRepo:  appRepo,
	}
}

func (u *permUsecase) Create(i *CreateInput) error {

	if i.AppID == "" {
		return errors.New("application ID is required")
	}

	if i.Code == "" {
		return errors.New("resource is required")
	}

	app, _ := u.appRepo.GetByID(i.AppID)
	if app == nil {
		return errors.New("application not found")
	}

	existing, _ := u.permRepo.GetByCode(i.Code)
	if existing != nil {
		return errors.New("permission already exists")
	}

	newPerm := &entity.Permission{
		ID:            idgen.NewUUIDv7(),
		ApplicationID: app.ID,
		Code:          i.Code,
		Description:   i.Description,
		Version:       1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return u.permRepo.Create(newPerm)
}

func (u *permUsecase) SyncPermissions(i *SyncInput) error {
	if i.AppCode == "" {
		return errors.New("application code is required")
	}

	app, _ := u.appRepo.GetByCode(i.AppCode)

	if app == nil {
		return errors.New("application not found")
	}

	for _, p := range i.Permissions {

		perm := &entity.Permission{
			ID:            idgen.NewUUIDv7(),
			ApplicationID: app.ID,
			Code:          p.Code,
			Description:   p.Description,
			Version:       1,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		err := u.permRepo.Upsert(perm)
		if err != nil {
			msg := fmt.Sprintf("sync permissios error %s", err)
			return errors.New(msg)
		}
	}

	return nil
}
