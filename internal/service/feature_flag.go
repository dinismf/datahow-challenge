package service

import (
	"context"
	"datahow-challenge/internal/domain"
	"errors"
)

type FeatureFlagService struct {
	featureFlagRepository  domain.IFeatureFlagRepository
	userOverrideRepository domain.IUserOverrideRepository
}

func NewFeatureFlagService(r domain.IFeatureFlagRepository, or domain.IUserOverrideRepository) *FeatureFlagService {
	return &FeatureFlagService{featureFlagRepository: r, userOverrideRepository: or}
}

func (s *FeatureFlagService) Create(ctx context.Context, req domain.CreateFeatureFlagRequest) (domain.FeatureFlagResponse, *domain.ServiceError) {
	flag := domain.NewFeatureFlag(req.Key, req.Name, req.GlobalEnabled)

	result, err := s.featureFlagRepository.Create(ctx, flag)
	if err != nil {
		if errors.Is(err, domain.ErrConflict) {
			return domain.FeatureFlagResponse{}, domain.ErrSvcConflict.WithReason(err)
		}
		return domain.FeatureFlagResponse{}, domain.ErrSvcInternal.WithReason(err)
	}

	return domain.NewFeatureFlagResponse(result), nil
}

func (s *FeatureFlagService) Get(ctx context.Context, key string) (domain.FeatureFlagResponse, *domain.ServiceError) {
	result, err := s.featureFlagRepository.GetByID(ctx, key)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.FeatureFlagResponse{}, domain.ErrSvcNotFound.WithReason(err)
		}
		return domain.FeatureFlagResponse{}, domain.ErrSvcInternal.WithReason(err)
	}

	return domain.NewFeatureFlagResponse(result), nil
}

func (s *FeatureFlagService) UpdateGlobal(ctx context.Context, key string, req domain.UpdateGlobalRequest) *domain.ServiceError {
	flag, err := s.featureFlagRepository.GetByID(ctx, key)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrSvcNotFound.WithReason(err)
		}
		return domain.ErrSvcInternal.WithReason(err)
	}

	flag.SetEnabled(req.Enabled)

	if _, err = s.featureFlagRepository.Update(ctx, flag); err != nil {
		return domain.ErrSvcInternal.WithReason(err)
	}

	return nil
}

// UpdateUserOverride upserts a per-user override for the given flag.
// The flag must exist; the override is created or replaced atomically.
func (s *FeatureFlagService) UpdateUserOverride(ctx context.Context, key, userID string, req domain.UpdateUserOverrideRequest) (domain.UserOverrideResponse, *domain.ServiceError) {
	if _, err := s.featureFlagRepository.GetByID(ctx, key); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.UserOverrideResponse{}, domain.ErrSvcNotFound.WithReason(err)
		}
		return domain.UserOverrideResponse{}, domain.ErrSvcInternal.WithReason(err)
	}

	result, err := s.userOverrideRepository.Set(ctx, domain.UserOverride{FlagId: key, UserId: userID, Enabled: req.Enabled})
	if err != nil {
		return domain.UserOverrideResponse{}, domain.ErrSvcInternal.WithReason(err)
	}

	return domain.NewUserOverrideResponse(result), nil
}

// EvaluateForUser resolves the effective state of a flag for a specific user.
// Evaluation is asymmetric:
//   - global=on                → always enabled; overrides are ignored
//   - global=off, override=on  → enabled  (user has early access before rollout)
//   - global=off, override=off → disabled (no override, or override matches global)
//   - global=off, no override  → disabled (global default)
func (s *FeatureFlagService) EvaluateForUser(ctx context.Context, key, userID string) (domain.EvaluationResponse, *domain.ServiceError) {
	flag, err := s.featureFlagRepository.GetByID(ctx, key)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.EvaluationResponse{}, domain.ErrSvcNotFound.WithReason(err)
		}
		return domain.EvaluationResponse{}, domain.ErrSvcInternal.WithReason(err)
	}

	if flag.IsEnabledGlobally() {
		return domain.NewEvaluationResponse(true, domain.EvaluationReasonGlobal), nil
	}

	override, err := s.userOverrideRepository.Get(ctx, key, userID)
	if err == nil {
		return domain.NewEvaluationResponse(override.Enabled, domain.EvaluationReasonUserOverride), nil
	}
	if !errors.Is(err, domain.ErrNotFound) {
		// ErrNotFound is expected (no override set); anything else is a storage failure.
		return domain.EvaluationResponse{}, domain.ErrSvcInternal.WithReason(err)
	}

	return domain.NewEvaluationResponse(false, domain.EvaluationReasonGlobal), nil
}
