package memory

import (
	"context"
	"datahow-challenge/internal/core"
	"sync"
	"time"
)

type FeatureFlagRepository struct {
	data map[string]core.FeatureFlag
	mu   sync.RWMutex
}

func NewFeatureFlagRepository() *FeatureFlagRepository {
	return &FeatureFlagRepository{
		data: make(map[string]core.FeatureFlag),
	}
}

func (r *FeatureFlagRepository) Create(_ context.Context, flag core.FeatureFlag) (core.FeatureFlag, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[flag.Id]; ok {
		return core.FeatureFlag{}, ErrKeyConflict
	}

	now := time.Now()
	flag.CreatedAt = now
	flag.UpdatedAt = now
	r.data[flag.Id] = flag
	return flag, nil
}

func (r *FeatureFlagRepository) GetByID(_ context.Context, id string) (core.FeatureFlag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	flag, ok := r.data[id]
	if !ok {
		return core.FeatureFlag{}, ErrKeyNotFound
	}
	return flag, nil
}

func (r *FeatureFlagRepository) Update(_ context.Context, flag core.FeatureFlag) (core.FeatureFlag, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.data[flag.Id]
	if !ok {
		return core.FeatureFlag{}, ErrKeyNotFound
	}

	flag.CreatedAt = existing.CreatedAt
	flag.UpdatedAt = time.Now()
	r.data[flag.Id] = flag
	return flag, nil
}

func (r *FeatureFlagRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.data, id)
	return nil
}
