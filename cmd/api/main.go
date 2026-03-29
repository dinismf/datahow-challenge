package main

import (
	"datahow-challenge/internal/infrastructure/memory"
	"datahow-challenge/internal/service"
	httpTransport "datahow-challenge/internal/transport/rest"

	"github.com/labstack/echo/v4"
)

func main() {

	repo := memory.NewFeatureFlagRepository()
	overrideRepo := memory.NewUserOverrideRepository()

	svc := service.NewFeatureFlagService(repo, overrideRepo)
	handler := httpTransport.NewHandler(svc)

	e := echo.New()
	httpTransport.RegisterRoutes(e, handler)

	e.Start(":8080")
}
