package main

import (
	"SSO/internal/app"
	"SSO/internal/config"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustConfig()
	fmt.Println(cfg)

	log := SetupLogger(cfg.Env)

	log.Info("Hello World")

	application := app.NewApp(log, cfg.Port, "ss", cfg.Timeout)

	go application.MustStart()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

	<-stop

	application.GracefulStop()
	log.Info("Goodbye")
	//TODO: инит app
	//TODO: Grpc
}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		log = slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		))
	}

	return log
}
