package auth

import (
	"context"
	ssov1 "github.com/DavidJackso/protos/gen/sso"
	"google.golang.org/grpc"
)

type serverApi struct {
	ssov1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverApi{})
}

func (s *serverApi) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	return &ssov1.LoginResponse{
		Token: "Hi",
	}, nil
}

func (s *serverApi) Register(
	ctx context.Context, req *ssov1.RegisterRequest) (
	*ssov1.RegisterResponse, error) {
	panic("implement me")
}
