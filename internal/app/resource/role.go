package resource

import (
	"github.com/jackc/pgx/v5/pgxpool"
	roleHdl "localdev.me/authorizer/internal/delivery/http/v1"
	"localdev.me/authorizer/internal/domain/repository"
	roleRepo "localdev.me/authorizer/internal/infrastructure/persistence/postgres/repository"
	roleUC "localdev.me/authorizer/internal/usecase/role"
)

type Role struct {
	Repository repository.RoleRepository
	Usecase    roleUC.Usecase
	Handler    *roleHdl.RoleHandler
}

func NewRole(db *pgxpool.Pool) *Role {
	repo := roleRepo.NewRoleRepositoryPGX(db)
	uc := roleUC.NewRoleUsecase(repo)
	hdl := roleHdl.NewRoleHandler(uc)

	return &Role{
		Repository: repo,
		Usecase:    uc,
		Handler:    hdl,
	}
}
