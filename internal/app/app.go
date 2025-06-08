package app

import (
	grpcapp "SSO/internal/app/grpc"
	"SSO/internal/config"
	authservice "SSO/internal/services/auth"
	"SSO/internal/storage/pg"
	"time"

	"log"
	"log/slog"
)

type App struct {
	gRPCServer *grpcapp.App
	log        *slog.Logger
}

func NewApp(log *slog.Logger, tokenTTL time.Duration, db config.DBConfig, grpcConf config.GRPCConfig) *App {
	storage, err := pg.New(log, db)
	if err != nil {
		log.Error("Failed to create storage", "error", err)
	}

	authService := authservice.New(
		log,
		storage,
		storage,
		storage,
		tokenTTL,
	)

	grpcApp := grpcapp.NewApp(log, authService, grpcConf.Port)

	return &App{
		gRPCServer: grpcApp,
		log:        log,
	}
}

func (a *App) MustStart() {
	err := a.gRPCServer.Run()
	if err != nil {
		panic(err)
	}
}

func (a *App) GracefulStop() {
	err := a.gRPCServer.Stop()
	if err != nil {
		log.Fatal(err)
	}
	a.log.Info("GRPC server gracefully stopped")
}
