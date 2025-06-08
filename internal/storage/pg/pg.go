package pg

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
)

type Storage struct {
	db *sql.DB
}

func New(log *slog.Logger, storage string) (*Storage, error) {
	const op = "storage.New"

	connStr := "user=postgres password=password dbname=productdb sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error("Unable to connect to database", "error", err, "storage", storage)
		return nil, fmt.Errorf(op+" %w", err)
	}
	return &Storage{
		db: db,
	}, nil
}

func (storage *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	smtm, err := storage.db.Prepare("INSERT INTO user (email, password_hash) VALUES (?, ?)")
}
