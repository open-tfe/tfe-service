package router

import (
	"github.com/gorilla/mux"
	"github.com/open-tfe/tfe-service/internal/api/handlers"
)

func (r *Router) registerOrganizationRoutes(api *mux.Router) {
	orgHandler := handlers.NewOrganizationHandler(r.service, r.logger)

	// Organizations endpoints
	api.HandleFunc("/organizations", orgHandler.List).Methods("GET")
	api.HandleFunc("/organizations", orgHandler.Create).Methods("POST")
	api.HandleFunc("/organizations/{name}", orgHandler.Read).Methods("GET")
	api.HandleFunc("/organizations/{name}", orgHandler.Update).Methods("PATCH")
	api.HandleFunc("/organizations/{name}", orgHandler.Delete).Methods("DELETE")

	// Organization entitlement set
	api.HandleFunc("/organizations/{name}/entitlement-set", orgHandler.ReadEntitlements).Methods("GET")

	// Organization relationships
	api.HandleFunc("/organizations/{name}/relationships/module-producers", orgHandler.ShowModuleProducers).Methods("GET")
	api.HandleFunc("/organizations/{name}/relationships/data-retention-policy", orgHandler.ShowDataRetentionPolicy).Methods("GET")
	api.HandleFunc("/organizations/{name}/relationships/data-retention-policy", orgHandler.UpdateDataRetentionPolicy).Methods("POST")
	api.HandleFunc("/organizations/{name}/relationships/data-retention-policy", orgHandler.DeleteDataRetentionPolicy).Methods("DELETE")
}
