package domain

import "errors"

type CreateFeatureFlagRequest struct {
	Key           string `json:"key"`
	Name          string `json:"name"`
	GlobalEnabled bool   `json:"global_enabled"`
}

func (r CreateFeatureFlagRequest) IsValid() error {
	if r.Key == "" {
		return errors.New("key is required")
	}
	if len(r.Key) > 255 {
		return errors.New("key exceeds maximum length of 255 characters")
	}
	if r.Name == "" {
		return errors.New("name is required")
	}
	if len(r.Name) > 255 {
		return errors.New("name exceeds maximum length of 255 characters")
	}

	return nil
}

type UpdateGlobalRequest struct {
	Enabled bool `json:"enabled"`
}

type UpdateUserOverrideRequest struct {
	Enabled bool `json:"enabled"`
}
