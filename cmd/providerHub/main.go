package main

import (
	"log/slog"
	"os"
	"os/signal"
	"providerHub/internal/app"
	"providerHub/internal/config"
	"providerHub/internal/lib/logger"
	"syscall"
)

// TODO: Написать свой генератор случайных чисел (math/rand)
// TODO: render.JSON
// TODO: Сделать удобный регистратор зависимостей

// TODO: Сервисы по работе с правами и пользователями
// TODO: Возможно сервис по работе с компонентами приложения по api keys?
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
