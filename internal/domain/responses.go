package domain

import "time"

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewErrorResponse(code string, message string) ErrorResponse {
	return ErrorResponse{Code: code, Message: message}
}

func NewBadRequestError(message string) ErrorResponse {
	return NewErrorResponse(ErrSvcInvalidInput.Code, message)
}

type FeatureFlagResponse struct {
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	GlobalEnabled bool      `json:"global_enabled"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewFeatureFlagResponse(flag FeatureFlag) FeatureFlagResponse {
	return FeatureFlagResponse{
		Id:            flag.Id,
		Name:          flag.Name,
		GlobalEnabled: flag.GlobalEnabled,
		CreatedAt:     flag.CreatedAt,
		UpdatedAt:     flag.UpdatedAt,
	}
}

type UserOverrideResponse struct {
	FlagId    string    `json:"flag_id"`
	UserId    string    `json:"user_id"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUserOverrideResponse(o UserOverride) UserOverrideResponse {
	return UserOverrideResponse{
		FlagId:    o.FlagId,
		UserId:    o.UserId,
		Enabled:   o.Enabled,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

const (
	EvaluationReasonUserOverride = "user_override"
	EvaluationReasonGlobal       = "global"
)

type EvaluationResponse struct {
	Enabled bool   `json:"enabled"`
	Reason  string `json:"reason,omitempty"`
}

func NewEvaluationResponse(enabled bool, reason string) EvaluationResponse {
	return EvaluationResponse{Enabled: enabled, Reason: reason}
}
