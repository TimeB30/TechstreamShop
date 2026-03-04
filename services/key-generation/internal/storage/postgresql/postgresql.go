package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgresql.New"
	db, err := pgxpool.New(context.Background(), storagePath)
	if err != nil {
		return nil, fmt.Errorf("db connection error %s %w", op, err)
	}
	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) SaveKey(userID int64, key string, startDate time.Time, endDate time.Time) (int64, error) {
	const op = "storage.postgresql.SaveKey"
	var id int64
	err := s.db.QueryRow(context.Background(), "INSERT INTO keys (user_id, start_date, end_date, key) VALUES ($1, $2, $3, $4) RETURNING ID", userID, startDate, endDate, key).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s %w", op, err)
	}

	return id, nil
}
