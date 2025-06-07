package auth

import (
	"context"
	ssov1 "github.com/DavidJackso/protos/gen/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string) (userID int64, err error)
}

type serverApi struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverApi{auth: auth})
}

// Login TODO: переписать валидацию
func (s *serverApi) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid email or password")
	}

	token, err := s.auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		//TODO дописать проверку ошибки
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil

	//TODO: implement login via service

	return nil, nil
}

func (s *serverApi) Register(
	ctx context.Context, req *ssov1.RegisterRequest) (
	*ssov1.RegisterResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid email or password")
	}
	userId, err := s.auth.RegisterNewUser(ctx, req.Email, req.Password)
	if err != nil {
		//TODO: тоже пересмотреть
		return nil, status.Error(codes.Internal, "failed to register")
	}
	return &ssov1.RegisterResponse{
		UserId: userId,
	}, nil
}
