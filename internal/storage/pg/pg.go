package pg

import (
	"SSO/internal/config"
	"SSO/internal/domain/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	"log/slog"
)

type Storage struct {
	db *sql.DB
}

func New(log *slog.Logger, config config.DBConfig) (*Storage, error) {
	const op = "storage.New"

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.User, config.Password, config.Database)

	log.Info("Connecting to PostgreSQL", op, connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error("Unable to connect to database", "error", err)
		return nil, fmt.Errorf(op+" %w", err)
	}

	// Проверяем соединение
	if err = db.Ping(); err != nil {
		log.Error("Cannot ping database", "error", err)
		return nil, fmt.Errorf(op+" %w", err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (storage *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.SaveUser"

	stmt, err := storage.db.PrepareContext(ctx, `INSERT INTO "users" (email, password_hash) VALUES ($1, $2) RETURNING id`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRowContext(ctx, email, passHash).Scan(&id)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// Код 23505 - уникальное ограничение нарушено (user уже существует)
			return 0, fmt.Errorf("%s: %w", op, errors.New("ErrUserExists"))
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (storage *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.User"

	stmt, err := storage.db.PrepareContext(ctx, `SELECT id, email, password_hash FROM "users" WHERE email=$1`)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var user models.User
	err = stmt.QueryRowContext(ctx, email).Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, errors.New("user not found"))
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (storage *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "storage.App"

	stmt, err := storage.db.PrepareContext(ctx, `SELECT id FROM apps WHERE id=$1`)
	if err != nil {
		log.Fatal(err)
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var app models.App
	err = stmt.QueryRowContext(ctx, appID).Scan(&app.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, errors.New("app not found"))
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
