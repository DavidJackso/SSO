package auth

import (
	"SSO/internal/domain/models"
	"SSO/internal/services/lib/jwt"
	"SSO/internal/storage"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (user models.User, err error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

// New returns new instance of the Auth
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		userSaver:    userSaver,
		userProvider: userProvider,
		log:          log,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (auth *Auth) Login(
	ctx context.Context,
	email, password string, appID int,
) (string, error) {
	op := "auth.Login"

	log := auth.log.With(slog.String("op", op),
		slog.String("email", email))

	log.Info("attempting to login user")

	user, err := auth.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			auth.log.Warn("invalid credentials")

			return " ", fmt.Errorf("invalid credentials: %w", err)

		}

		auth.log.Error("failed to get user")

		return " ", fmt.Errorf("failed to get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		auth.log.Warn("invalid password")

		return " ", fmt.Errorf("invalid password: %w", err)
	}

	app, err := auth.appProvider.App(ctx, appID)
	log.Info("successfully logged in")

	token, err := jwt.NewToken(user, app, auth.tokenTTL)
	if err != nil {
		auth.log.Error("failed to create token")

		return " ", fmt.Errorf("failed to create token: %w", err)
	}
	return token, nil
}

func (auth *Auth) RegisterNewUser(
	ctx context.Context,
	email, password string,
) (int64, error) {
	op := "auth.RegisterNewUser"

	log := auth.log.With(
		slog.String("op", op),
		slog.String("email", email))

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password")

		return 0, fmt.Errorf("failed to generate password: %w", err)
	}

	id, err := auth.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			auth.log.Warn("user already exists")
			return 0, fmt.Errorf("user already exists: %w", err)
		}
		log.Error("failed to save user")
		return 0, fmt.Errorf("failed to save user: %w", err)
	}

	log.Info("successfully saved user")

	return id, nil
}
