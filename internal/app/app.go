package app

import (
	"log/slog"
	"os"

	httpApp "github.com/AlonMell/auth/internal/app/http"
	"github.com/AlonMell/auth/internal/app/postgres"
	"github.com/AlonMell/auth/internal/app/router"
	"github.com/AlonMell/auth/internal/infra/lib/jwt"
	"github.com/AlonMell/auth/internal/infra/repo"
	"github.com/AlonMell/auth/internal/service/auth"
	"github.com/AlonMell/auth/internal/service/user"
	"github.com/AlonMell/grovelog/util"
	sq "github.com/Masterminds/squirrel"
)

var (
	postgresPlaceholder = sq.Dollar
)

type App struct {
	Server *httpApp.Server
}

func New(
	log *slog.Logger,
	postgresCfg postgres.Config,
	jwtCfg jwt.Config,
	serverCfg httpApp.Config,
) *App {
	db, err := postgres.New(postgresCfg, log)
	if err != nil {
		log.Error("error with start db postgres!", util.Err(err))
		os.Exit(1)
	}

	userRepo := repo.NewUserRepo(db, postgresPlaceholder)

	authService := auth.New(log, userRepo, jwtCfg)
	userService := user.New(log, userRepo)

	mux := router.New(log, authService, userService)
	mux.Prepare(jwtCfg)

	server := httpApp.New(log, serverCfg, mux)

	return &App{Server: server}
}
