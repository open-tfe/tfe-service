package service

import (
	"context"

	"github.com/google/uuid"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/open-tfe/tfe-service/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	// Organization methods
	ListOrganizations(ctx context.Context, query string) ([]*tfe.Organization, error)
	CreateOrganization(ctx context.Context, org *tfe.Organization) (*tfe.Organization, error)
	ReadOrganization(ctx context.Context, name string) (*tfe.Organization, error)
	UpdateOrganization(ctx context.Context, name string, org *tfe.Organization) error
	DeleteOrganization(ctx context.Context, name string) error
	GetOrganizationIDByName(ctx context.Context, name string) (uuid.UUID, error)
	ReadOrganizationEntitlements(ctx context.Context, name string) (*tfe.Entitlements, error)

	// Project methods
	ListProjects(ctx context.Context, orgID uuid.UUID) ([]*models.Project, []*tfe.Project, error)
	CreateProject(ctx context.Context, project *tfe.Project) (*tfe.Project, error)
	ReadProject(ctx context.Context, projectID string) (*tfe.Project, error)
	UpdateProject(ctx context.Context, project *tfe.Project) (*tfe.Project, error)
	DeleteProject(ctx context.Context, projectID string) error

	// User methods
	ListUsers(ctx context.Context) ([]*tfe.User, error)
	CreateUser(ctx context.Context, user *tfe.User) (*tfe.User, error)
	ReadUser(ctx context.Context, userID string) (*tfe.User, error)
	UpdateUser(ctx context.Context, userID string, user *tfe.User) (*tfe.User, error)
	DeleteUser(ctx context.Context, userID string) error
	ReadCurrentUser(ctx context.Context) (*tfe.User, error)
}

type service struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewService(db *gorm.DB, logger *zap.Logger) Service {
	return &service{
		db:     db,
		logger: logger.With(zap.String("component", "services")),
	}
}
