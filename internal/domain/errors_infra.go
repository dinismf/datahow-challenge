package domain

import "errors"

// Sentinel errors returned by infrastructure layers.
// Repositories wrap storage-specific errors (e.g. pgconn.PgError, redis.Nil)
// into these so the service layer stays decoupled from any storage backend.
var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrInvalidInput = errors.New("invalid input")
	ErrInternal     = errors.New("internal error")
)
