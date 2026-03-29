package memory

import (
	"context"
	"datahow-challenge/internal/core"
	"sync"
	"time"
)

type UserOverrideRepository struct {
	data map[string]core.UserOverride // key: flagId:userId
	mu   sync.RWMutex
}

func NewUserOverrideRepository() *UserOverrideRepository {
	return &UserOverrideRepository{
		data: make(map[string]core.UserOverride),
	}
}

func (r *UserOverrideRepository) Set(_ context.Context, override core.UserOverride) (core.UserOverride, error) {
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

func (r *UserOverrideRepository) Get(_ context.Context, flagId, userId string) (core.UserOverride, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := flagId + ":" + userId
	o, ok := r.data[key]
	if !ok {
		return core.UserOverride{}, ErrKeyNotFound
	}
	return o, nil
}

func (r *UserOverrideRepository) Delete(_ context.Context, flagId, userId string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.data, flagId+":"+userId)
	return nil
}
