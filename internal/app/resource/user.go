package resource

import (
	"github.com/jackc/pgx/v5/pgxpool"
	hdl "github.com/mafzaidi/authorizer/internal/delivery/http/v1"
	"github.com/mafzaidi/authorizer/internal/domain/repository"
	repo "github.com/mafzaidi/authorizer/internal/infrastructure/persistence/postgres/repository"
	uc "github.com/mafzaidi/authorizer/internal/usecase/user"
)

type User struct {
	UserRepository     repository.UserRepository
	RoleRepostory      repository.RoleRepository
	UserRoleRepository repository.UserRoleRepository
	Usecase            uc.Usecase
	Handler            *hdl.UserHandler
}

func NewUser(db *pgxpool.Pool) *User {
	userRepo := repo.NewUserRepositoryPGX(db)
	roleRepo := repo.NewRoleRepositoryPGX(db)
	userRoleRepo := repo.NewUserRoleRepositoryPGX(db)
	userUC := uc.NewUserUsecase(userRepo, roleRepo, userRoleRepo)
	userHdl := hdl.NewUserHandler(userUC)

	return &User{
		UserRepository:     userRepo,
		RoleRepostory:      roleRepo,
		UserRoleRepository: userRoleRepo,
		Usecase:            userUC,
		Handler:            userHdl,
	}
}
