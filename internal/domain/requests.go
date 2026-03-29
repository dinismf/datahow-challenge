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
	if r.Name == "" {
		return errors.New("name is required")
	}

	return nil
}

type UpdateGlobalRequest struct {
	Enabled bool `json:"enabled"`
}

type UpdateUserOverrideRequest struct {
	Enabled bool `json:"enabled"`
}
