package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/open-tfe/tfe-service/internal/auth"
	"github.com/open-tfe/tfe-service/internal/constants"
	"github.com/open-tfe/tfe-service/internal/service"
	"go.uber.org/zap"
)

// Router holds all the route configurations
type Router struct {
	*mux.Router
	service service.Service
	logger  *zap.Logger
}

// NewRouter creates and configures a new router
func NewRouter(jwtSecret string, service service.Service, logger *zap.Logger) *Router {
	r := &Router{
		Router:  mux.NewRouter(),
		service: service,
		logger:  logger,
	}

	// API v2 routes
	api := r.PathPrefix(constants.APIVersionPath).Subrouter()
	api.Use(auth.JWTMiddleware(jwtSecret, r.logger))

	// Register all routes
	r.registerOrganizationRoutes(api)
	r.registerProjectRoutes(api)
	r.registerUserRoutes(api)

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		logger.Debug("Not Found",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
		)
	})

	return r
}
