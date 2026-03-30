package memory

import (
	"datahow-challenge/internal/domain"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserOverrideRepository_Set(t *testing.T) {
	t.Run("creates override and sets timestamps", func(t *testing.T) {
		repo := NewUserOverrideRepository()
		o := domain.UserOverride{FlagId: "my-flag", UserId: "user-1", Enabled: true}

		result, err := repo.Set(ctx, o)

		require.NoError(t, err)
		assert.Equal(t, "my-flag", result.FlagId)
		assert.Equal(t, "user-1", result.UserId)
		assert.True(t, result.Enabled)
		assert.False(t, result.CreatedAt.IsZero())
		assert.False(t, result.UpdatedAt.IsZero())
	})

	t.Run("upsert preserves CreatedAt on subsequent calls", func(t *testing.T) {
		repo := NewUserOverrideRepository()
		o := domain.UserOverride{FlagId: "my-flag", UserId: "user-1", Enabled: true}
		first, _ := repo.Set(ctx, o)

		o.Enabled = false
		second, err := repo.Set(ctx, o)

		require.NoError(t, err)
		assert.False(t, second.Enabled)
		assert.Equal(t, first.CreatedAt, second.CreatedAt)
		assert.True(t, second.UpdatedAt.Equal(first.UpdatedAt) || second.UpdatedAt.After(first.UpdatedAt))
	})
}

func TestUserOverrideRepository_Get(t *testing.T) {
	t.Run("returns stored override", func(t *testing.T) {
		repo := NewUserOverrideRepository()
		o := domain.UserOverride{FlagId: "my-flag", UserId: "user-1", Enabled: true}
		repo.Set(ctx, o)

		result, err := repo.Get(ctx, "my-flag", "user-1")

		require.NoError(t, err)
		assert.Equal(t, "my-flag", result.FlagId)
		assert.Equal(t, "user-1", result.UserId)
		assert.True(t, result.Enabled)
	})

	t.Run("returns ErrInfraNotFound for missing key", func(t *testing.T) {
		repo := NewUserOverrideRepository()

		_, err := repo.Get(ctx, "my-flag", "user-1")

		require.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInfraNotFound))
	})

	t.Run("overrides are scoped per flag and user", func(t *testing.T) {
		repo := NewUserOverrideRepository()
		repo.Set(ctx, domain.UserOverride{FlagId: "flag-a", UserId: "user-1", Enabled: true})

		_, err := repo.Get(ctx, "flag-b", "user-1")

		assert.True(t, errors.Is(err, domain.ErrInfraNotFound))
	})
}
