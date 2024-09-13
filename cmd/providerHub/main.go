package main

import (
	"github.com/AlonMell/ProviderHub/cmd/providerHub/config"
	loader "github.com/AlonMell/ProviderHub/internal/infra/lib/config"
	"github.com/AlonMell/ProviderHub/internal/infra/lib/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlonMell/ProviderHub/internal/app"
)

// TODO: Написать свой генератор случайных чисел (math/rand)
// TODO: Сделать удобный регистратор зависимостей

// TODO: Сервисы по работе с правами и пользователями
// TODO: Возможно сервис по работе с компонентами приложения по api keys?

// @title ProviderHub API
// @version 1.0
// @description This is a sample server ProviderHub server.

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	var cfg config.Config
	loader.MustLoad(&cfg)

	log := logger.SetupLogger(cfg.Env)

	log.Info("starting server", slog.Any("cfg", cfg))
	log.Debug("debug messages are enabled")

	application := app.New(log, cfg.Postgres, cfg.JWT, cfg.HTTPServer)

	go application.Server.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	signalReceived := <-stop

	log.Info("stopping application", slog.String("signal", signalReceived.String()))

	application.Server.Stop()

	log.Info("server stopped")
}
