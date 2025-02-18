package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/jsonapi"
	"github.com/open-tfe/tfe-service/internal/service"
	"go.uber.org/zap"
)

type UserHandler struct {
	svc    service.Service
	logger *zap.Logger
}

func NewUserHandler(svc service.Service, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		svc:    svc,
		logger: logger.With(zap.String("handler", "user")),
	}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.ListUsers(r.Context())
	if err != nil {
		h.logger.Error("failed to list users", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	err = jsonapi.MarshalPayload(w, users)
	if err != nil {
		h.logger.Error("failed to marshal response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user tfe.User
	if err := jsonapi.UnmarshalPayload(r.Body, &user); err != nil {
		h.logger.Error("failed to decode request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.svc.CreateUser(r.Context(), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(http.StatusCreated)
	err = jsonapi.MarshalPayload(w, createdUser)
	if err != nil {
		h.logger.Error("failed to marshal response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) Read(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	user, err := h.svc.ReadUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	err = jsonapi.MarshalPayload(w, user)
	if err != nil {
		h.logger.Error("failed to marshal response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	var user tfe.User
	if err := jsonapi.UnmarshalPayload(r.Body, &user); err != nil {
		h.logger.Error("failed to decode request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedUser, err := h.svc.UpdateUser(r.Context(), userID, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	err = jsonapi.MarshalPayload(w, updatedUser)
	if err != nil {
		h.logger.Error("failed to marshal response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if err := h.svc.DeleteUser(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) AccountDetails(w http.ResponseWriter, r *http.Request) {
	user, err := h.svc.ReadCurrentUser(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	err = jsonapi.MarshalPayload(w, user)
	if err != nil {
		h.logger.Error("failed to marshal response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
