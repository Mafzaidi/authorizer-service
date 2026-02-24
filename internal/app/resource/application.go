package resource

import (
	"github.com/jackc/pgx/v5/pgxpool"
	appHdl "github.com/mafzaidi/authorizer/internal/delivery/http/v1"
	"github.com/mafzaidi/authorizer/internal/domain/repository"
	appRepo "github.com/mafzaidi/authorizer/internal/infrastructure/persistence/postgres/repository"
	appUC "github.com/mafzaidi/authorizer/internal/usecase/application"
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
