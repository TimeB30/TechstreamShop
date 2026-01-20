package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	db *pgx.Conn
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgresql.New"
	db, err := pgx.Connect(context.Background(), storagePath)
	if err != nil {
		return nil, fmt.Errorf("db connection error %s", op)
	}
	return &Storage{
		db: db,
	}, nil
}
