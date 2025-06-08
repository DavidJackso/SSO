package app

import (
	grpcapp "SSO/internal/app/grpc"
	"SSO/internal/services/auth"
	"log"
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

	authService := auth.New(log)

	grpcApp := grpcapp.NewApp(
		log, auth, grpcPort)
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
