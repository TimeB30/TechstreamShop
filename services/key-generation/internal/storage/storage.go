package storage

import "errors"

var (
	ErrURLNotFound = errors.New("key not found")
	ErrURLExists   = errors.New("key exists")
)
