package rest

import "github.com/labstack/echo/v4"

func RegisterRoutes(e *echo.Echo, h *Handler) {
	e.POST("/flags", h.Create)
	e.GET("/flags/:id", h.Get)
	e.PUT("/flags/:id/global", h.UpdateGlobal)
	e.PUT("/flags/:id/users/:user_id", h.UpdateUserOverride)
	e.GET("/flags/:id/users/:user_id/evaluation", h.EvaluateUser)

}
