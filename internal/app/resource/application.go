package resource

import (
	"github.com/jackc/pgx/v5/pgxpool"
	appHdl "localdev.me/authorizer/internal/delivery/http/v1"
	"localdev.me/authorizer/internal/domain/repository"
	appRepo "localdev.me/authorizer/internal/infrastructure/persistence/postgres/repository"
	appUC "localdev.me/authorizer/internal/usecase/application"
)

type App struct {
	Repository repository.AppRepository
	Usecase    appUC.Usecase
	Handler    *appHdl.AppHandler
}

func NewApp(db *pgxpool.Pool) *App {
	repo := appRepo.NewAppRepositoryPGX(db)
	uc := appUC.NewAppUsecase(repo)
	hdl := appHdl.NewAppHandler(uc)

	return &App{
		Repository: repo,
		Usecase:    uc,
		Handler:    hdl,
	}
}
