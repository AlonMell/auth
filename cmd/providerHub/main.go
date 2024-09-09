package main

import (
	"log/slog"
	"os"
	"os/signal"
	"providerHub/internal/infra/config"
	"providerHub/internal/infra/lib/logger"
	"syscall"

	"providerHub/internal/app"
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

	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("starting server", slog.Any("cfg", cfg))
	log.Debug("debug messages are enabled")

	application := app.New(log, cfg)

	go application.Server.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	signalReceived := <-stop

	log.Info("stopping application", slog.String("signal", signalReceived.String()))

	application.Server.Stop()

	log.Info("server stopped")
}
