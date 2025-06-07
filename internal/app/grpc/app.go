package grpcapp

import (
	"SSO/internal/grpc/auth"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"strconv"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewApp(log *slog.Logger, port int) *App {
	GRPCServer := grpc.NewServer()

	auth.Register(GRPCServer)

	return &App{
		log:        log,
		gRPCServer: GRPCServer,
		port:       port,
	}
}

func (a *App) Run() error {
	const op = "grpcApp.Run"

	log := a.log.With(
		slog.String("op", "op"),
		slog.String("port", strconv.Itoa(a.port)))

	l, err := net.Listen("tcp", ":"+strconv.Itoa(a.port))
	if err != nil {
		return fmt.Errorf(op+" failed to listen: %w", err)
	}

	log.Info("starting gRPC server", "port", a.port, l.Addr().String())

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf(op+" failed to serve: %w", err)
	}

	return nil
}

func (a *App) Stop() error {
	const op = "grpcApp.Stop"

	log := a.log.With(slog.String("op", "op"),
		slog.String("port", strconv.Itoa(a.port)))

	a.gRPCServer.GracefulStop()

	log.Info("stopped gRPC server")

	return nil
}
