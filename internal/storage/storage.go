package storage

import "errors"

var (
	ErrEventNotFound = errors.New("event not found")
	ErrEventExists   = errors.New("event already exists")
)
