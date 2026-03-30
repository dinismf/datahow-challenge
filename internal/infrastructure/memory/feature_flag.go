package memory

import (
	"context"
	"datahow-challenge/internal/domain"
	"fmt"
	"sync"
	"time"
)

type FeatureFlagInMemoryRepository struct {
	name string
	data map[string]domain.FeatureFlag
	mu   sync.RWMutex
}

func NewFeatureFlagRepository() *FeatureFlagInMemoryRepository {
	return &FeatureFlagInMemoryRepository{
		name: "memory.FeatureFlagInMemoryRepository",
		data: make(map[string]domain.FeatureFlag),
	}
}

func (r *FeatureFlagInMemoryRepository) Create(_ context.Context, flag domain.FeatureFlag) (domain.FeatureFlag, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[flag.Id]; ok {
		return domain.FeatureFlag{}, fmt.Errorf("%s: key %q: %w", r.name, flag.Id, domain.ErrInfraConflict)
	}

	now := time.Now()
	flag.CreatedAt = now
	flag.UpdatedAt = now
	r.data[flag.Id] = flag
	return flag, nil
}

func (r *FeatureFlagInMemoryRepository) GetByID(_ context.Context, id string) (domain.FeatureFlag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	flag, ok := r.data[id]
	if !ok {
		return domain.FeatureFlag{}, fmt.Errorf("%s: key %q: %w", r.name, id, domain.ErrInfraNotFound)
	}
	return flag, nil
}

func (r *FeatureFlagInMemoryRepository) Update(_ context.Context, flag domain.FeatureFlag) (domain.FeatureFlag, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.data[flag.Id]
	if !ok {
		return domain.FeatureFlag{}, fmt.Errorf("%s: key %q: %w", r.name, flag.Id, domain.ErrInfraNotFound)
	}

	flag.CreatedAt = existing.CreatedAt
	flag.UpdatedAt = time.Now()
	r.data[flag.Id] = flag
	return flag, nil
}

func (r *FeatureFlagInMemoryRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.data, id)
	return nil
}
