package memory

import (
	"context"
	"datahow-challenge/internal/domain"
	"fmt"
	"sync"
	"time"
)

type UserOverrideInMemoryRepository struct {
	name string
	data map[string]domain.UserOverride // key: flagId:userId
	mu   sync.RWMutex
}

func NewUserOverrideRepository() *UserOverrideInMemoryRepository {
	return &UserOverrideInMemoryRepository{
		name: "memory.UserOverrideInMemoryRepository",
		data: make(map[string]domain.UserOverride),
	}
}

func (r *UserOverrideInMemoryRepository) Set(_ context.Context, override domain.UserOverride) (domain.UserOverride, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := override.FlagId + ":" + override.UserId
	now := time.Now()
	if existing, ok := r.data[key]; ok {
		override.CreatedAt = existing.CreatedAt
	} else {
		override.CreatedAt = now
	}
	override.UpdatedAt = now
	r.data[key] = override
	return override, nil
}

func (r *UserOverrideInMemoryRepository) Get(_ context.Context, flagId, userId string) (domain.UserOverride, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := flagId + ":" + userId
	o, ok := r.data[key]
	if !ok {
		return domain.UserOverride{}, fmt.Errorf("%s: key %q: %w", r.name, key, domain.ErrInfraNotFound)
	}
	return o, nil
}
