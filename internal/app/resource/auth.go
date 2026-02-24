package resource

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mafzaidi/authorizer/config"
	authHdl "github.com/mafzaidi/authorizer/internal/delivery/http/v1"
	"github.com/mafzaidi/authorizer/internal/domain/repository"
	authRepoPGX "github.com/mafzaidi/authorizer/internal/infrastructure/persistence/postgres/repository"
	authRepoRedis "github.com/mafzaidi/authorizer/internal/infrastructure/persistence/redis/repository"
	authSvc "github.com/mafzaidi/authorizer/internal/service"
	authUC "github.com/mafzaidi/authorizer/internal/usecase/auth"
	"github.com/redis/go-redis/v9"
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
