package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/open-tfe/tfe-service/internal/models"
	"github.com/open-tfe/tfe-service/internal/service"
	"go.uber.org/zap"
)

type OrganizationHandler struct {
	svc    service.OrganizationService
	logger *zap.Logger
}

func NewOrganizationHandler(svc service.OrganizationService, logger *zap.Logger) *OrganizationHandler {
	return &OrganizationHandler{
		svc:    svc,
		logger: logger.With(zap.String("handler", "organization")),
	}
}

// @Summary List organizations
// @Description Get all organizations
// @Tags organizations
// @Accept json
// @Produce json
// @Param query query string false "Search query"
// @Success 200 {array} models.Organization
// @Router /organizations [get]
func (h *OrganizationHandler) List(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("listing organizations", zap.String("query", r.URL.Query().Get("q")))

	query := r.URL.Query().Get("q")
	orgs, err := h.svc.List(r.Context(), query)
	if err != nil {
		h.logger.Error("failed to list organizations", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to JSON API format
	data := make([]interface{}, 0, len(orgs))
	for _, org := range orgs {
		data = append(data, map[string]interface{}{
			"id":         org.Name,
			"type":       "organizations",
			"attributes": org,
		})
	}

	resp := map[string]interface{}{
		"data": data,
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary Create organization
// @Description Create a new organization
// @Tags organizations
// @Accept json
// @Produce json
// @Param organization body models.Organization true "Organization object"
// @Success 201 {object} models.Organization
// @Router /organizations [post]
func (h *OrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("creating organization")

	var org models.Organization
	if err := json.NewDecoder(r.Body).Decode(&org); err != nil {
		h.logger.Error("failed to decode request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdOrg, err := h.svc.Create(r.Context(), org.ToTFE())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdOrg)
}

// @Summary Read organization
// @Description Get an organization by name
// @Tags organizations
// @Accept json
// @Produce json
// @Param name path string true "Organization name"
// @Success 200 {object} models.Organization
// @Router /organizations/{name} [get]
func (h *OrganizationHandler) Read(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	org, err := h.svc.Read(r.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := map[string]interface{}{
		"data": map[string]interface{}{
			"id":         org.Name,
			"type":       "organizations",
			"attributes": org,
		},
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary Update organization
// @Description Update an organization by name
// @Tags organizations
// @Accept json
// @Produce json
// @Param name path string true "Organization name"
// @Param organization body models.Organization true "Organization object"
// @Success 200 {object} models.Organization
// @Router /organizations/{name} [patch]
func (h *OrganizationHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	var org models.Organization
	if err := json.NewDecoder(r.Body).Decode(&org); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.svc.Update(r.Context(), name, org.ToTFE()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(org)
}

// @Summary Delete organization
// @Description Delete an organization by name
// @Tags organizations
// @Accept json
// @Produce json
// @Param name path string true "Organization name"
// @Success 204 "No Content"
// @Router /organizations/{name} [delete]
func (h *OrganizationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if err := h.svc.Delete(r.Context(), name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *OrganizationHandler) ShowEntitlementSet(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement entitlement set logic
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": map[string]interface{}{
			"id":   "org-entitlement",
			"type": "entitlement-set",
			// Add other entitlement attributes here
		},
	})
}

func (h *OrganizationHandler) ShowModuleProducers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": []interface{}{}})
}

func (h *OrganizationHandler) ShowDataRetentionPolicy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": nil})
}

func (h *OrganizationHandler) UpdateDataRetentionPolicy(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *OrganizationHandler) DeleteDataRetentionPolicy(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
