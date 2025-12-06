package resource

import (
	"github.com/jackc/pgx/v5/pgxpool"
	hdl "localdev.me/authorizer/internal/delivery/http/v1"
	"localdev.me/authorizer/internal/domain/repository"
	repo "localdev.me/authorizer/internal/infrastructure/persistence/postgres/repository"
	uc "localdev.me/authorizer/internal/usecase/permission"
)

type Perm struct {
	PermRepository repository.PermRepository
	AppRepository  repository.AppRepository
	Usecase        uc.Usecase
	Handler        *hdl.PermHandler
}

func NePerm(db *pgxpool.Pool) *Perm {
	permRepo := repo.NewPermRepositoryPGX(db)
	appRepo := repo.NewAppRepositoryPGX(db)
	uc := uc.NewPermUsecase(permRepo, appRepo)
	hdl := hdl.NewPermHandler(uc)

	return &Perm{
		PermRepository: permRepo,
		AppRepository:  appRepo,
		Usecase:        uc,
		Handler:        hdl,
	}
}
