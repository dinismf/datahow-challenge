package rest

import (
	"context"
	"datahow-challenge/internal/domain"

	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Create(_ context.Context, flag domain.FeatureFlag) (domain.FeatureFlag, error) {
	args := m.Called(flag)
	return args.Get(0).(domain.FeatureFlag), args.Error(1)
}

func (m *mockRepository) GetByID(_ context.Context, id string) (domain.FeatureFlag, error) {
	args := m.Called(id)
	return args.Get(0).(domain.FeatureFlag), args.Error(1)
}

func (m *mockRepository) Update(_ context.Context, flag domain.FeatureFlag) (domain.FeatureFlag, error) {
	args := m.Called(flag)
	return args.Get(0).(domain.FeatureFlag), args.Error(1)
}

func (m *mockRepository) Delete(_ context.Context, id string) error {
	return m.Called(id).Error(0)
}

type mockOverrideRepository struct {
	mock.Mock
}

func (m *mockOverrideRepository) Set(_ context.Context, o domain.UserOverride) (domain.UserOverride, error) {
	args := m.Called(o)
	return args.Get(0).(domain.UserOverride), args.Error(1)
}

func (m *mockOverrideRepository) Get(_ context.Context, flagId, userId string) (domain.UserOverride, error) {
	args := m.Called(flagId, userId)
	return args.Get(0).(domain.UserOverride), args.Error(1)
}

func (m *mockOverrideRepository) Delete(_ context.Context, flagId, userId string) error {
	return m.Called(flagId, userId).Error(0)
}
