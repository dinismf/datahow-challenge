package memory

import (
	"datahow-challenge/internal/core"
	"fmt"
)

// Storage-specific errors for the in-memory backend.
// Each wraps the corresponding core sentinel so that callers using
// errors.Is(err, core.ErrNotFound) still work, while also giving
// diagnostics that are specific to this storage implementation.
//
// As the system grows, other backends (Postgres, Redis, …) define their own
// equivalent files (e.g. postgres.ErrRowNotFound) that wrap the same core
// sentinels. The service layer never sees the backend-specific type.
var (
	ErrKeyNotFound = fmt.Errorf("key not found in store: %w", core.ErrNotFound)
	ErrKeyConflict = fmt.Errorf("key already exists in store: %w", core.ErrConflict)
)
