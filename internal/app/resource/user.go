package resource

import (
	"github.com/jackc/pgx/v5/pgxpool"
	userHdl "localdev.me/authorizer/internal/delivery/http/v1"
	"localdev.me/authorizer/internal/domain/repository"
	userRepo "localdev.me/authorizer/internal/infrastructure/persistence/postgres/repository"
	userUC "localdev.me/authorizer/internal/usecase/user"
)

type User struct {
	Repository repository.UserRepository
	Usecase    userUC.Usecase
	Handler    *userHdl.UserHandler
}

func NewUser(db *pgxpool.Pool) *User {
	repo := userRepo.NewUserRepositoryPGX(db)
	uc := userUC.NewUserUsecase(repo)
	hdl := userHdl.NewUserHandler(uc)

	return &User{
		Repository: repo,
		Usecase:    uc,
		Handler:    hdl,
	}
}
