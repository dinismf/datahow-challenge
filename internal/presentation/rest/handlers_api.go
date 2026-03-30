package rest

import (
	"datahow-challenge/internal/domain"
	"datahow-challenge/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *service.FeatureFlagService
}

func NewHandler(s *service.FeatureFlagService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Create(c echo.Context) error {

	ctx := c.Request().Context()

	var req domain.CreateFeatureFlagRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.NewBadRequestError(err.Error()))
	}

	if err := req.IsValid(); err != nil {
		return c.JSON(http.StatusBadRequest, domain.NewBadRequestError(err.Error()))
	}

	created, svcErr := h.service.Create(ctx, req)
	if svcErr != nil {
		return httpError(c, svcErr)
	}

	return c.JSON(http.StatusCreated, created)
}

func (h *Handler) Get(c echo.Context) error {

	ctx := c.Request().Context()

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, domain.NewBadRequestError("id is required"))
	}

	flag, svcErr := h.service.Get(ctx, id)
	if svcErr != nil {
		return httpError(c, svcErr)
	}

	return c.JSON(http.StatusOK, flag)
}

func (h *Handler) UpdateGlobal(c echo.Context) error {

	ctx := c.Request().Context()

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, domain.NewBadRequestError("id is required"))
	}

	var req domain.UpdateGlobalRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.NewBadRequestError(err.Error()))
	}

	if svcErr := h.service.UpdateGlobal(ctx, id, req); svcErr != nil {
		return httpError(c, svcErr)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) UpdateUserOverride(c echo.Context) error {

	ctx := c.Request().Context()

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, domain.NewBadRequestError("id is required"))
	}

	userID := c.Param("user_id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, domain.NewBadRequestError("user_id is required"))
	}

	var req domain.UpdateUserOverrideRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.NewBadRequestError(err.Error()))
	}

	override, svcErr := h.service.UpdateUserOverride(ctx, id, userID, req)
	if svcErr != nil {
		return httpError(c, svcErr)
	}

	return c.JSON(http.StatusOK, override)
}

func (h *Handler) EvaluateUser(c echo.Context) error {

	ctx := c.Request().Context()

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, domain.NewBadRequestError("id is required"))
	}

	userID := c.Param("user_id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, domain.NewBadRequestError("user_id is required"))
	}

	response, svcErr := h.service.EvaluateForUser(ctx, id, userID)
	if svcErr != nil {
		return httpError(c, svcErr)
	}

	return c.JSON(http.StatusOK, response)
}
