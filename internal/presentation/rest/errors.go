package rest

import (
	"datahow-challenge/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

func httpError(c echo.Context, err *domain.ServiceError) error {
	// Log the full chain (Message + Reason) server-side so operators have
	// complete context. Only the sanitised Message reaches the client below.
	c.Logger().Error(err.LogError())

	switch err.Code {
	case domain.ErrSvcNotFound.Code:
		return c.JSON(http.StatusNotFound, domain.ErrorResponse{Code: err.Code, Message: err.Message})
	case domain.ErrSvcConflict.Code:
		return c.JSON(http.StatusConflict, domain.ErrorResponse{Code: err.Code, Message: err.Message})
	case domain.ErrSvcInvalidInput.Code:
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{Code: err.Code, Message: err.Message})
	default:
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Code: domain.ErrSvcInternal.Code, Message: "internal server error"})
	}
}
