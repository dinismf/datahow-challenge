package domain

import (
	"context"
	"time"
)

type IFeatureFlagRepository interface {
	Create(ctx context.Context, flag FeatureFlag) (FeatureFlag, error)
	GetByID(ctx context.Context, id string) (FeatureFlag, error)
	Update(ctx context.Context, flag FeatureFlag) (FeatureFlag, error)
	Delete(ctx context.Context, id string) error
}

type FeatureFlag struct {
	Id            string
	Name          string
	GlobalEnabled bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewFeatureFlag(key string, name string, enabled bool) FeatureFlag {
	return FeatureFlag{
		Id:            key,
		Name:          name,
		GlobalEnabled: enabled,
	}
}

func (f *FeatureFlag) IsEnabledGlobally() bool {
	return f.GlobalEnabled
}

func (f *FeatureFlag) SetEnabled(enabled bool) {
	f.GlobalEnabled = enabled
}
