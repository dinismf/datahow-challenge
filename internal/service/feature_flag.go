package service

import (
	"context"
	"datahow-challenge/internal/core"
	"errors"
)

type FeatureFlagService struct {
	featureFlagRepository  core.IFeatureFlagRepository
	userOverrideRepository core.IUserOverrideRepository
}

func NewFeatureFlagService(r core.IFeatureFlagRepository, or core.IUserOverrideRepository) *FeatureFlagService {
	return &FeatureFlagService{featureFlagRepository: r, userOverrideRepository: or}
}

func (s *FeatureFlagService) Create(ctx context.Context, req core.CreateFeatureFlagRequest) (core.FeatureFlagResponse, *core.ServiceError) {
	flag := core.NewFeatureFlag(req.Key, req.Name, req.GlobalEnabled)
	result, err := s.featureFlagRepository.Create(ctx, flag)
	if err != nil {
		if errors.Is(err, core.ErrConflict) {
			return core.FeatureFlagResponse{}, core.ErrSvcConflict.WithReason(err)
		}
		return core.FeatureFlagResponse{}, core.ErrSvcInternal.WithReason(err)
	}
	return core.NewFeatureFlagResponse(result), nil
}

func (s *FeatureFlagService) Get(ctx context.Context, key string) (core.FeatureFlagResponse, *core.ServiceError) {
	result, err := s.featureFlagRepository.GetByID(ctx, key)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return core.FeatureFlagResponse{}, core.ErrSvcNotFound.WithReason(err)
		}
		return core.FeatureFlagResponse{}, core.ErrSvcInternal.WithReason(err)
	}
	return core.NewFeatureFlagResponse(result), nil
}

func (s *FeatureFlagService) UpdateGlobal(ctx context.Context, key string, req core.UpdateGlobalRequest) *core.ServiceError {
	flag, err := s.featureFlagRepository.GetByID(ctx, key)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return core.ErrSvcNotFound.WithReason(err)
		}
		return core.ErrSvcInternal.WithReason(err)
	}
	flag.SetEnabled(req.Enabled)
	if _, err = s.featureFlagRepository.Update(ctx, flag); err != nil {
		return core.ErrSvcInternal.WithReason(err)
	}
	return nil
}

// UpdateUserOverride upserts a per-user override for the given flag.
// The flag must exist; the override is created or replaced atomically.
func (s *FeatureFlagService) UpdateUserOverride(ctx context.Context, key, userID string, req core.UpdateUserOverrideRequest) (core.UserOverrideResponse, *core.ServiceError) {
	if _, err := s.featureFlagRepository.GetByID(ctx, key); err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return core.UserOverrideResponse{}, core.ErrSvcNotFound.WithReason(err)
		}
		return core.UserOverrideResponse{}, core.ErrSvcInternal.WithReason(err)
	}
	result, err := s.userOverrideRepository.Set(ctx, core.UserOverride{FlagId: key, UserId: userID, Enabled: req.Enabled})
	if err != nil {
		return core.UserOverrideResponse{}, core.ErrSvcInternal.WithReason(err)
	}
	return core.NewUserOverrideResponse(result), nil
}

// EvaluateForUser resolves the effective state of a flag for a specific user.
// A user-level override always wins over the global setting, in both directions:
//   - global=off, override=on  → enabled  (user opted in before feature rollout)
//   - global=on,  override=off → disabled (user opted out of active rollout)
//
// If no override exists, the global setting is returned.
func (s *FeatureFlagService) EvaluateForUser(ctx context.Context, key, userID string) (core.EvaluationResponse, *core.ServiceError) {
	flag, err := s.featureFlagRepository.GetByID(ctx, key)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return core.EvaluationResponse{}, core.ErrSvcNotFound.WithReason(err)
		}
		return core.EvaluationResponse{}, core.ErrSvcInternal.WithReason(err)
	}

	override, err := s.userOverrideRepository.Get(ctx, key, userID)
	if err == nil {
		return core.NewEvaluationResponse(override.Enabled, core.EvaluationReasonUserOverride), nil
	}
	if !errors.Is(err, core.ErrNotFound) {
		// ErrNotFound is expected (no override set); anything else is a storage failure.
		return core.EvaluationResponse{}, core.ErrSvcInternal.WithReason(err)
	}

	return core.NewEvaluationResponse(flag.IsEnabledGlobally(), core.EvaluationReasonGlobal), nil
}
