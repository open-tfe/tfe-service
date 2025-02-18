package router

import (
	"github.com/gorilla/mux"
	"github.com/open-tfe/tfe-service/internal/api/handlers"
)

func (r *Router) registerUserRoutes(api *mux.Router) {
	userHandler := handlers.NewUserHandler(r.service, r.logger)

	// User endpoints
	api.HandleFunc("/users", userHandler.List).Methods("GET")
	api.HandleFunc("/users", userHandler.Create).Methods("POST")
	api.HandleFunc("/users/{user_id}", userHandler.Read).Methods("GET")
	api.HandleFunc("/users/{user_id}", userHandler.Update).Methods("PATCH")
	api.HandleFunc("/users/{user_id}", userHandler.Delete).Methods("DELETE")
	api.HandleFunc("/account/details", userHandler.AccountDetails).Methods("GET")
}
