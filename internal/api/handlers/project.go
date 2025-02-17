package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/jsonapi"
	"github.com/open-tfe/tfe-service/internal/service"
	"go.uber.org/zap"
)

type ProjectHandler struct {
	projectService      service.ProjectService
	organizationService service.OrganizationService
	logger              *zap.Logger
}

func NewProjectHandler(projectService service.ProjectService, organizationService service.OrganizationService, logger *zap.Logger) *ProjectHandler {
	return &ProjectHandler{
		projectService:      projectService,
		organizationService: organizationService,
		logger:              logger.With(zap.String("handler", "project")),
	}
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgName := vars["organization_name"]
	h.logger.Debug("listing projects", zap.String("organization_name", orgName))

	orgID, err := h.organizationService.GetIDByName(r.Context(), orgName)
	if err != nil {
		h.logger.Error("failed to get organization ID", zap.String("organization_name", orgName), zap.Error(err))
		http.Error(w, "Organization not found", http.StatusNotFound)
		return
	}
	h.logger.Debug("Get organization ID", zap.String("organization_id", orgID.String()))
	projects, err := h.projectService.List(r.Context(), orgID)
	if err != nil {
		h.logger.Error("failed to list projects", zap.String("organization_id", orgID.String()), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	err = jsonapi.MarshalPayload(w, projects)
	if err != nil {
		h.logger.Error("failed to marshal response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("creating project")

	var project *tfe.Project
	if err := jsonapi.UnmarshalPayload(r.Body, &project); err != nil {
		h.logger.Error("failed to decode request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdProject, err := h.projectService.Create(r.Context(), project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(http.StatusCreated)
	err = jsonapi.MarshalPayload(w, createdProject)
	if err != nil {
		h.logger.Error("failed to marshal response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *ProjectHandler) Read(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["project_id"]

	project, err := h.projectService.Read(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	err = jsonapi.MarshalPayload(w, project)
	if err != nil {
		h.logger.Error("failed to marshal response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["project_id"]

	var project *tfe.Project
	if err := jsonapi.UnmarshalPayload(r.Body, &project); err != nil {
		h.logger.Error("failed to decode request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	project.ID = projectID

	updatedProject, err := h.projectService.Update(r.Context(), project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/vnd.api+json")
	err = jsonapi.MarshalPayload(w, updatedProject)
	if err != nil {
		h.logger.Error("failed to marshal response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["project_id"]

	if err := h.projectService.Delete(r.Context(), projectID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
