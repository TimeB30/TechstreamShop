package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/timeb30/techstreamshop/services/telegram-bot/internal/storage"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgresql.new"
	db, err := pgxpool.New(context.Background(), storagePath)
	if err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) AddUser(ctx context.Context, page *storage.Page) error {
	const op = "storage.postgresql.addUser"
	_, err := s.db.Exec(ctx, `INSERT INTO users (id, first_name, last_name, username) VALUES ($1, $2, $3, $4)`,
		page.UserId,
		page.UserFirstName,
		page.UserLastName,
		page.UserName,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetUser(ctx context.Context, userId int64) (*storage.Page, error) {
	const op = "storage.postgresql.getUser"
	res := storage.Page{}
	err := s.db.QueryRow(ctx, `SELECT * FROM users WHERE id = $1`, userId).Scan(&res.UserId, &res.UserFirstName, &res.UserLastName, &res.UserName)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
