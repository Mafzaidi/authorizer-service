package resource

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"localdev.me/authorizer/config"
	authHdl "localdev.me/authorizer/internal/delivery/http/v1"
	"localdev.me/authorizer/internal/domain/repository"
	authRepoPGX "localdev.me/authorizer/internal/infrastructure/persistence/postgres/repository"
	authRepoRedis "localdev.me/authorizer/internal/infrastructure/persistence/redis/repository"
	authSvc "localdev.me/authorizer/internal/service"
	authUC "localdev.me/authorizer/internal/usecase/auth"
)

type Auth struct {
	AuthRepository     repository.AuthRepository
	UserRepository     repository.UserRepository
	AppRepository      repository.AppRepository
	UserRoleRepository repository.UserRoleRepository
	RoleRepository     repository.RoleRepository
	RolePermRepository repository.RolePermRepository
	PermRepository     repository.PermRepository
	JWTService         authSvc.JWTService
	Usecase            authUC.Usecase
	Handler            *authHdl.AuthHandler
}

func NewAuth(db *pgxpool.Pool, cfg *config.Config, redis *redis.Client) *Auth {
	authRepo := authRepoRedis.NewAuthRepository(redis)
	userRepo := authRepoPGX.NewUserRepositoryPGX(db)
	appRepo := authRepoPGX.NewAppRepositoryPGX(db)
	userRoleRepo := authRepoPGX.NewUserRoleRepositoryPGX(db)
	roleRepo := authRepoPGX.NewRoleRepositoryPGX(db)
	rolePermRepo := authRepoPGX.NewRolePermRepositoryPGX(db)
	permRepo := authRepoPGX.NewPermRepositoryPGX(db)
	jwtSvc := authSvc.NewJWTService(authRepo, userRepo, appRepo, userRoleRepo, roleRepo, rolePermRepo, permRepo)
	uc := authUC.NewAuthUseCase(authRepo, userRepo, jwtSvc)
	hdl := authHdl.NewAuthHandler(uc, cfg)

	return &Auth{
		AuthRepository:     authRepo,
		UserRepository:     userRepo,
		AppRepository:      appRepo,
		UserRoleRepository: userRoleRepo,
		RoleRepository:     roleRepo,
		PermRepository:     permRepo,
		Usecase:            uc,
		Handler:            hdl,
	}
}
