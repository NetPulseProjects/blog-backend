package app

import (
	"app/internal/config"
	"app/internal/repository"
	"app/internal/service"
	"app/internal/transport/http"
	"app/internal/transport/http/handler"
	"app/pkg/infra/database/postgresql"
	"app/pkg/infra/logger/handlers/slogpretty"
	"app/pkg/infra/logger/sl"
	"log/slog"
	"os"
)

type App struct {
	Log *slog.Logger
	Cfg *config.Config
}

func (app *App) New() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	app.Cfg = cfg

	log.Info("Starting application", slog.String("env", cfg.Env))

	// Dependencies
	sqlStorage, err := postgresql.New(postgresql.CreateConnectionString(app.Cfg))
	if err != nil {
		app.Log.Error("Failed to init repository", sl.Err(err))
		os.Exit(1)
	}

	// Repos
	repositories := repository.NewRepositories(sqlStorage)

	// Services
	services := service.NewServices(service.Deps{
		Repos:  repositories,
		Config: cfg,
	})

	// HTTP Handler
	handlers := handler.NewTransportHandler(log, cfg, services)

	// HTTP Server
	http.NewTransportServer(
		log,
		cfg,
		handlers,
	)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case config.EnvLocal:
		{
			opts := slogpretty.PrettyHandlerOptions{
				SlogOpts: &slog.HandlerOptions{
					Level: slog.LevelDebug,
				},
			}

			loggerHandler := opts.NewPrettyHandler(os.Stdout)
			log = slog.New(loggerHandler)
		}

	case config.EnvDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case config.EnvProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
