package domain

import "errors"

// Contract errors for the infrastructure layer.
// Repositories wrap storage-specific errors (e.g. pgconn.PgError, redis.Nil)
// into these via %w so the service layer stays decoupled from any storage backend.
var (
	ErrInfraNotFound     = errors.New("not found")
	ErrInfraConflict     = errors.New("conflict")
	ErrInfraInvalidInput = errors.New("invalid input")
	ErrInfraInternal     = errors.New("internal error")
)
