package storage

import "context"

type Storage interface {
	AddUser(ctx context.Context, page *Page) error
	GetUser(ctx context.Context, userId int64) (*Page, error)
}

type Page struct {
	UserId        int64
	UserFirstName string
	UserLastName  string
	UserName      string
}
