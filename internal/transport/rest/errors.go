package rest

import (
	"datahow-challenge/internal/core"
	"net/http"

	"github.com/labstack/echo/v4"
)

func httpError(c echo.Context, err *core.ServiceError) error {
	// Log the full chain (Message + Reason) server-side so operators have
	// complete context. Only the sanitised Message reaches the client below.
	c.Logger().Error(err.LogError())

	switch err.Code {
	case core.ErrSvcNotFound.Code:
		return c.JSON(http.StatusNotFound, core.ErrorResponse{Code: err.Code, Message: err.Message})
	case core.ErrSvcConflict.Code:
		return c.JSON(http.StatusConflict, core.ErrorResponse{Code: err.Code, Message: err.Message})
	case core.ErrSvcInvalidInput.Code:
		return c.JSON(http.StatusBadRequest, core.ErrorResponse{Code: err.Code, Message: err.Message})
	default:
		return c.JSON(http.StatusInternalServerError, core.ErrorResponse{Code: core.ErrSvcInternal.Code, Message: "internal server error"})
	}
}
