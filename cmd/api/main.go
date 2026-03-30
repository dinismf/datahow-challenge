package main

import (
	"datahow-challenge/internal/infrastructure/memory"
	httpTransport "datahow-challenge/internal/presentation/rest"
	"datahow-challenge/internal/service"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize in-memory repositories for feature flags and user overrides
	// Using in-memory storage for simplicity; can be swapped with persistent storage (Postgres, Redis) by implementing domain interfaces
	featureFlagRepo := memory.NewFeatureFlagRepository()
	overrideRepo := memory.NewUserOverrideRepository()

	// Create service layer with injected repositories
	// Service encapsulates business logic and coordinates between repositories
	svc := service.NewFeatureFlagService(featureFlagRepo, overrideRepo)

	// Create HTTP handler with injected service
	// Handler translates HTTP requests/responses and delegates to service layer
	handler := httpTransport.NewHandler(svc)

	// Initialize Echo web framework
	e := echo.New()

	// Register all REST API routes with the handler
	// Routes include CRUD operations for flags, global toggles, user overrides, and evaluation
	httpTransport.RegisterRoutes(e, handler)

	// Start HTTP server on port 8080
	e.Start(":8080")
}
