package resource

import (
	"github.com/jackc/pgx/v5/pgxpool"
	hdl "localdev.me/authorizer/internal/delivery/http/v1"
	"localdev.me/authorizer/internal/domain/repository"
	repo "localdev.me/authorizer/internal/infrastructure/persistence/postgres/repository"
	uc "localdev.me/authorizer/internal/usecase/role"
)

type Role struct {
	RoleRepository     repository.RoleRepository
	AppRepository      repository.AppRepository
	PermRepository     repository.PermRepository
	RolePermRepository repository.RolePermRepository
	Usecase            uc.Usecase
	Handler            *hdl.RoleHandler
}

func NewRole(db *pgxpool.Pool) *Role {
	roleRepo := repo.NewRoleRepositoryPGX(db)
	appRepo := repo.NewAppRepositoryPGX(db)
	permRepo := repo.NewPermRepositoryPGX(db)
	rolePermRepo := repo.NewRolePermRepositoryPGX(db)
	roleUC := uc.NewRoleUsecase(roleRepo, appRepo, permRepo, rolePermRepo)
	roleHdl := hdl.NewRoleHandler(roleUC)

	return &Role{
		RoleRepository:     roleRepo,
		AppRepository:      appRepo,
		PermRepository:     permRepo,
		RolePermRepository: rolePermRepo,
		Usecase:            roleUC,
		Handler:            roleHdl,
	}
}
