package resource

import (
	"github.com/jackc/pgx/v5/pgxpool"
	hdl "github.com/mafzaidi/authorizer/internal/delivery/http/v1"
	"github.com/mafzaidi/authorizer/internal/domain/repository"
	repo "github.com/mafzaidi/authorizer/internal/infrastructure/persistence/postgres/repository"
	uc "github.com/mafzaidi/authorizer/internal/usecase/permission"
)

type Perm struct {
	PermRepository repository.PermRepository
	AppRepository  repository.AppRepository
	Usecase        uc.Usecase
	Handler        *hdl.PermHandler
}

func NewPerm(db *pgxpool.Pool) *Perm {
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
