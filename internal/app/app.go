package app

import (
	grpcapp "SSO/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	gRPCServer *grpcapp.App
	log        *slog.Logger
}

func NewApp(log *slog.Logger, grpcPort int, storage string, tokenTTL time.Duration) *App {
	//TODO:db
	//TODO: init auth service

	grpcApp := grpcapp.NewApp(
		log, grpcPort)
	return &App{
		gRPCServer: grpcApp,
		log:        log,
	}
}

func (a *App) Start() error {
	err := a.gRPCServer.Run()
	if err != nil {
		a.log.Debug("Failed to start application", "error", err)
	}
	return err
}
