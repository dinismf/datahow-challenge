package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"datahow-challenge/internal/domain"
	"datahow-challenge/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stack wires mock repos → real service → real handler → real router.
// Tests assert on HTTP responses, covering the full request/response contract.
func stack(t *testing.T) (*mockRepository, *mockOverrideRepository, *echo.Echo) {
	t.Helper()
	repo := &mockRepository{}
	overrides := &mockOverrideRepository{}
	svc := service.NewFeatureFlagService(repo, overrides)
	e := echo.New()
	RegisterRoutes(e, NewHandler(svc))
	return repo, overrides, e
}

func do(e *echo.Echo, method, path, body string) *httptest.ResponseRecorder {
	var reqBody *strings.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	} else {
		reqBody = strings.NewReader("")
	}
	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func decodeBody[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var v T
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&v))
	return v
}

// --- POST /flags ---

func TestCreateHandler(t *testing.T) {
	t.Run("creates flag and returns 201", func(t *testing.T) {
		repo, _, e := stack(t)
		flag := domain.NewFeatureFlag("new-dashboard", "New Dashboard", false)
		repo.On("Create", flag).Return(flag, nil)

		rec := do(e, http.MethodPost, "/flags", `{"key":"new-dashboard","name":"New Dashboard","global_enabled":false}`)

		assert.Equal(t, http.StatusCreated, rec.Code)
		resp := decodeBody[domain.FeatureFlagResponse](t, rec)
		assert.Equal(t, "new-dashboard", resp.Id)
		assert.Equal(t, "New Dashboard", resp.Name)
		repo.AssertExpectations(t)
	})

	t.Run("returns 409 when key already exists", func(t *testing.T) {
		repo, _, e := stack(t)
		flag := domain.NewFeatureFlag("existing", "Existing", false)
		repo.On("Create", flag).Return(domain.FeatureFlag{}, domain.ErrInfraConflict)

		rec := do(e, http.MethodPost, "/flags", `{"key":"existing","name":"Existing"}`)

		assert.Equal(t, http.StatusConflict, rec.Code)
		repo.AssertExpectations(t)
	})

	t.Run("returns 400 when key is missing", func(t *testing.T) {
		_, _, e := stack(t)
		rec := do(e, http.MethodPost, "/flags", `{"name":"No Key"}`)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 400 when name is missing", func(t *testing.T) {
		_, _, e := stack(t)
		rec := do(e, http.MethodPost, "/flags", `{"key":"no-name"}`)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// --- GET /flags/:id ---

func TestGetHandler(t *testing.T) {
	t.Run("returns flag on success", func(t *testing.T) {
		repo, _, e := stack(t)
		repo.On("GetByID", "my-flag").Return(domain.NewFeatureFlag("my-flag", "My Flag", true), nil)

		rec := do(e, http.MethodGet, "/flags/my-flag", "")

		assert.Equal(t, http.StatusOK, rec.Code)
		resp := decodeBody[domain.FeatureFlagResponse](t, rec)
		assert.Equal(t, "my-flag", resp.Id)
		assert.True(t, resp.GlobalEnabled)
		repo.AssertExpectations(t)
	})

	t.Run("returns 404 when flag does not exist", func(t *testing.T) {
		repo, _, e := stack(t)
		repo.On("GetByID", "missing").Return(domain.FeatureFlag{}, domain.ErrInfraNotFound)

		rec := do(e, http.MethodGet, "/flags/missing", "")

		assert.Equal(t, http.StatusNotFound, rec.Code)
		repo.AssertExpectations(t)
	})
}

// --- PUT /flags/:id/global ---

func TestUpdateGlobalHandler(t *testing.T) {
	t.Run("returns 204 on success", func(t *testing.T) {
		repo, _, e := stack(t)
		flag := domain.NewFeatureFlag("my-flag", "My Flag", false)
		updated := domain.NewFeatureFlag("my-flag", "My Flag", true)
		repo.On("GetByID", "my-flag").Return(flag, nil)
		repo.On("Update", updated).Return(updated, nil)

		rec := do(e, http.MethodPut, "/flags/my-flag/global", `{"enabled":true}`)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		repo.AssertExpectations(t)
	})

	t.Run("returns 404 when flag does not exist", func(t *testing.T) {
		repo, _, e := stack(t)
		repo.On("GetByID", "missing").Return(domain.FeatureFlag{}, domain.ErrInfraNotFound)

		rec := do(e, http.MethodPut, "/flags/missing/global", `{"enabled":true}`)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		repo.AssertExpectations(t)
	})

	t.Run("returns 500 when update fails", func(t *testing.T) {
		repo, _, e := stack(t)
		flag := domain.NewFeatureFlag("my-flag", "My Flag", false)
		updated := domain.NewFeatureFlag("my-flag", "My Flag", true)
		repo.On("GetByID", "my-flag").Return(flag, nil)
		repo.On("Update", updated).Return(domain.FeatureFlag{}, domain.ErrInfraInternal)

		rec := do(e, http.MethodPut, "/flags/my-flag/global", `{"enabled":true}`)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		repo.AssertExpectations(t)
	})
}

// --- PUT /flags/:id/users/:user_id ---

func TestUpdateUserOverrideHandler(t *testing.T) {
	t.Run("returns override on success", func(t *testing.T) {
		repo, overrides, e := stack(t)
		repo.On("GetByID", "my-flag").Return(domain.NewFeatureFlag("my-flag", "My Flag", false), nil)
		o := domain.UserOverride{FlagId: "my-flag", UserId: "user-1", Enabled: true}
		overrides.On("Set", o).Return(o, nil)

		rec := do(e, http.MethodPut, "/flags/my-flag/users/user-1", `{"enabled":true}`)

		assert.Equal(t, http.StatusOK, rec.Code)
		resp := decodeBody[domain.UserOverrideResponse](t, rec)
		assert.Equal(t, "my-flag", resp.FlagId)
		assert.Equal(t, "user-1", resp.UserId)
		assert.True(t, resp.Enabled)
		repo.AssertExpectations(t)
		overrides.AssertExpectations(t)
	})

	t.Run("returns 404 when flag does not exist", func(t *testing.T) {
		repo, _, e := stack(t)
		repo.On("GetByID", "missing").Return(domain.FeatureFlag{}, domain.ErrInfraNotFound)

		rec := do(e, http.MethodPut, "/flags/missing/users/user-1", `{"enabled":true}`)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		repo.AssertExpectations(t)
	})
}

// --- GET /flags/:id/users/:user_id/evaluation ---

func TestEvaluateUserHandler(t *testing.T) {
	t.Run("returns global state when global is on", func(t *testing.T) {
		repo, overrides, e := stack(t)
		repo.On("GetByID", "my-flag").Return(domain.NewFeatureFlag("my-flag", "My Flag", true), nil)

		rec := do(e, http.MethodGet, "/flags/my-flag/users/user-1/evaluation", "")

		assert.Equal(t, http.StatusOK, rec.Code)
		resp := decodeBody[domain.EvaluationResponse](t, rec)
		assert.True(t, resp.Enabled)
		assert.Equal(t, domain.EvaluationReasonGlobal, resp.Reason)
		repo.AssertExpectations(t)
		overrides.AssertExpectations(t)
	})

	t.Run("returns global state when global is off and no override", func(t *testing.T) {
		repo, overrides, e := stack(t)
		repo.On("GetByID", "my-flag").Return(domain.NewFeatureFlag("my-flag", "My Flag", false), nil)
		overrides.On("Get", "my-flag", "user-1").Return(domain.UserOverride{}, domain.ErrInfraNotFound)

		rec := do(e, http.MethodGet, "/flags/my-flag/users/user-1/evaluation", "")

		assert.Equal(t, http.StatusOK, rec.Code)
		resp := decodeBody[domain.EvaluationResponse](t, rec)
		assert.False(t, resp.Enabled)
		assert.Equal(t, domain.EvaluationReasonGlobal, resp.Reason)
		repo.AssertExpectations(t)
		overrides.AssertExpectations(t)
	})

	t.Run("user override wins over global", func(t *testing.T) {
		repo, overrides, e := stack(t)
		repo.On("GetByID", "my-flag").Return(domain.NewFeatureFlag("my-flag", "My Flag", false), nil)
		overrides.On("Get", "my-flag", "user-1").Return(
			domain.UserOverride{FlagId: "my-flag", UserId: "user-1", Enabled: true}, nil,
		)

		rec := do(e, http.MethodGet, "/flags/my-flag/users/user-1/evaluation", "")

		assert.Equal(t, http.StatusOK, rec.Code)
		resp := decodeBody[domain.EvaluationResponse](t, rec)
		assert.True(t, resp.Enabled)
		assert.Equal(t, domain.EvaluationReasonUserOverride, resp.Reason)
		repo.AssertExpectations(t)
		overrides.AssertExpectations(t)
	})

	t.Run("global on wins even when override is off", func(t *testing.T) {
		repo, overrides, e := stack(t)
		repo.On("GetByID", "my-flag").Return(domain.NewFeatureFlag("my-flag", "My Flag", true), nil)

		rec := do(e, http.MethodGet, "/flags/my-flag/users/user-1/evaluation", "")

		assert.Equal(t, http.StatusOK, rec.Code)
		resp := decodeBody[domain.EvaluationResponse](t, rec)
		assert.True(t, resp.Enabled)
		assert.Equal(t, domain.EvaluationReasonGlobal, resp.Reason)
		repo.AssertExpectations(t)
		overrides.AssertExpectations(t)
	})

	t.Run("returns 404 when flag does not exist", func(t *testing.T) {
		repo, overrides, e := stack(t)
		repo.On("GetByID", "my-flag").Return(domain.FeatureFlag{}, domain.ErrInfraNotFound)

		rec := do(e, http.MethodGet, "/flags/my-flag/users/user-1/evaluation", "")

		assert.Equal(t, http.StatusNotFound, rec.Code)
		repo.AssertExpectations(t)
		overrides.AssertExpectations(t)
	})
}
