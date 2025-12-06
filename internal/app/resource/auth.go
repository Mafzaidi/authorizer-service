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
	AppRepository      repository.AppRepository
	UserRoleRepository repository.UserRoleRepository
	RoleRepository     repository.RoleRepository
	RolePermRepository repository.RolePermRepository
	PermRepository     repository.PermRepository
	Usecase            authUC.Usecase
	Handler            *authHdl.AuthHandler
}

func NewAuth(db *pgxpool.Pool, cfg *config.Config) *Auth {
	userRepo := authRepo.NewUserRepositoryPGX(db)
	appRepo := authRepo.NewAppRepositoryPGX(db)
	userRoleRepo := authRepo.NewUserRoleRepositoryPGX(db)
	roleRepo := authRepo.NewRoleRepositoryPGX(db)
	rolePermRepo := authRepo.NewRolePermRepositoryPGX(db)
	permRepo := authRepo.NewPermRepositoryPGX(db)
	uc := authUC.NewAuthUseCase(userRepo, appRepo, userRoleRepo, roleRepo, rolePermRepo, permRepo)
	hdl := authHdl.NewAuthHandler(uc, cfg)

	return &Auth{
		UserRepository:     userRepo,
		AppRepository:      appRepo,
		UserRoleRepository: userRoleRepo,
		RoleRepository:     roleRepo,
		PermRepository:     permRepo,
		Usecase:            uc,
		Handler:            hdl,
	}
}
