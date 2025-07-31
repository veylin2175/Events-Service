package storage

import "errors"

var (
	ErrURLNotFound = errors.New("event not found")
	ErrURLExists   = errors.New("event already exists")
)
