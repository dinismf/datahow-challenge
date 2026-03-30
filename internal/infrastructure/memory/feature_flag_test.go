package memory

import (
	"context"
	"datahow-challenge/internal/domain"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func TestFeatureFlagRepository_Create(t *testing.T) {
	t.Run("stores flag and sets timestamps", func(t *testing.T) {
		repo := NewFeatureFlagRepository()
		flag := domain.NewFeatureFlag("my-flag", "My Flag", true)

		result, err := repo.Create(ctx, flag)

		require.NoError(t, err)
		assert.Equal(t, "my-flag", result.Id)
		assert.False(t, result.CreatedAt.IsZero())
		assert.False(t, result.UpdatedAt.IsZero())
	})

	t.Run("returns ErrInfraConflict when key already exists", func(t *testing.T) {
		repo := NewFeatureFlagRepository()
		flag := domain.NewFeatureFlag("my-flag", "My Flag", false)
		_, err := repo.Create(ctx, flag)
		require.NoError(t, err)

		_, err = repo.Create(ctx, flag)

		require.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInfraConflict))
	})
}

func TestFeatureFlagRepository_GetByID(t *testing.T) {
	t.Run("returns stored flag", func(t *testing.T) {
		repo := NewFeatureFlagRepository()
		flag := domain.NewFeatureFlag("my-flag", "My Flag", true)
		created, _ := repo.Create(ctx, flag)

		result, err := repo.GetByID(ctx, "my-flag")

		require.NoError(t, err)
		assert.Equal(t, created, result)
	})

	t.Run("returns ErrInfraNotFound for missing key", func(t *testing.T) {
		repo := NewFeatureFlagRepository()

		_, err := repo.GetByID(ctx, "missing")

		require.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInfraNotFound))
	})
}

func TestFeatureFlagRepository_Update(t *testing.T) {
	t.Run("updates flag and preserves CreatedAt", func(t *testing.T) {
		repo := NewFeatureFlagRepository()
		created, _ := repo.Create(ctx, domain.NewFeatureFlag("my-flag", "My Flag", false))

		updated, err := repo.Update(ctx, domain.NewFeatureFlag("my-flag", "My Flag", true))

		require.NoError(t, err)
		assert.True(t, updated.GlobalEnabled)
		assert.Equal(t, created.CreatedAt, updated.CreatedAt)
		assert.True(t, updated.UpdatedAt.After(created.UpdatedAt) || updated.UpdatedAt.Equal(created.UpdatedAt))
	})

	t.Run("returns ErrInfraNotFound for missing key", func(t *testing.T) {
		repo := NewFeatureFlagRepository()

		_, err := repo.Update(ctx, domain.NewFeatureFlag("missing", "x", false))

		require.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInfraNotFound))
	})
}

func TestFeatureFlagRepository_Delete(t *testing.T) {
	t.Run("removes the flag", func(t *testing.T) {
		repo := NewFeatureFlagRepository()
		repo.Create(ctx, domain.NewFeatureFlag("my-flag", "My Flag", false))

		err := repo.Delete(ctx, "my-flag")
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, "my-flag")
		assert.True(t, errors.Is(err, domain.ErrInfraNotFound))
	})

	t.Run("is a no-op for missing key", func(t *testing.T) {
		repo := NewFeatureFlagRepository()

		err := repo.Delete(ctx, "missing")

		require.NoError(t, err)
	})
}
