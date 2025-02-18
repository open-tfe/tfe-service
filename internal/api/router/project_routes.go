package router

import (
	"github.com/gorilla/mux"
	"github.com/open-tfe/tfe-service/internal/api/handlers"
)

func (r *Router) registerProjectRoutes(api *mux.Router) {
	projectHandler := handlers.NewProjectHandler(r.service, r.logger)

	// Project endpoints
	api.HandleFunc("/organizations/{organization_name}/projects", projectHandler.List).Methods("GET")
	api.HandleFunc("/organizations/{organization_name}/projects", projectHandler.Create).Methods("POST")
	api.HandleFunc("/projects/{project_id}", projectHandler.Read).Methods("GET")
	api.HandleFunc("/projects/{project_id}", projectHandler.Update).Methods("PATCH")
	api.HandleFunc("/projects/{project_id}", projectHandler.Delete).Methods("DELETE")

	// Project tag bindings
	// api.HandleFunc("/projects/{project_id}/tag-bindings", projectHandler.ListTagBindings).Methods("GET")
	// api.HandleFunc("/projects/{project_id}/effective-tag-bindings", projectHandler.ListEffectiveTagBindings).Methods("GET")

	// Project workspace relationships
	// api.HandleFunc("/projects/{project_id}/relationships/workspaces", projectHandler.MoveWorkspaces).Methods("POST")
}
