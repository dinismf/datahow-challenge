package main

import (
	"datahow-challenge/internal/infrastructure/memory"
	httpTransport "datahow-challenge/internal/presentation/rest"
	"datahow-challenge/internal/service"

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
