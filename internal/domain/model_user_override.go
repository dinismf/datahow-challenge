package domain

import (
	"context"
	"time"
)

type IUserOverrideRepository interface {
	Set(ctx context.Context, override UserOverride) (UserOverride, error)
	Get(ctx context.Context, flagId, userId string) (UserOverride, error)
}

type UserOverride struct {
	FlagId    string
	UserId    string
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
