package router

import (
	"github.com/gorilla/mux"
	"github.com/open-tfe/tfe-service/internal/service"
	"go.uber.org/zap"
)

// Router holds all the route configurations
type Router struct {
	*mux.Router
	services *service.Services
	logger   *zap.Logger
}

// NewRouter creates and configures a new router
func NewRouter(services *service.Services, logger *zap.Logger) *Router {
	r := &Router{
		Router:   mux.NewRouter(),
		services: services,
		logger:   logger,
	}

	// API v2 routes
	api := r.PathPrefix("/api/v2").Subrouter()

	// Register all routes
	r.registerOrganizationRoutes(api)
	r.registerProjectRoutes(api)
	r.registerUserRoutes(api)

	return r
}
