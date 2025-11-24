package resource

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"localdev.me/authorizer/config"
	authHdl "localdev.me/authorizer/internal/delivery/http/v1"
	"localdev.me/authorizer/internal/domain/repository"
	authRepo "localdev.me/authorizer/internal/infrastructure/persistence/postgres/repository"
	authUC "localdev.me/authorizer/internal/usecase/auth"
)

type Auth struct {
	UserRepository     repository.UserRepository
	UserRoleRepository repository.UserRoleRepository
	RoleRepository     repository.RoleRepository
	Usecase            authUC.Usecase
	Handler            *authHdl.AuthHandler
}

func NewAuth(db *pgxpool.Pool, cfg *config.Config) *Auth {
	userRepo := authRepo.NewUserRepositoryPGX(db)
	userRoleRepo := authRepo.NewUserRoleRepositoryPGX(db)
	roleRepo := authRepo.NewRoleRepositoryPGX(db)
	uc := authUC.NewAuthUseCase(userRepo, userRoleRepo, roleRepo)
	hdl := authHdl.NewAuthHandler(uc, cfg)

	return &Auth{
		UserRepository:     userRepo,
		UserRoleRepository: userRoleRepo,
		RoleRepository:     roleRepo,
		Usecase:            uc,
		Handler:            hdl,
	}
}
