package service

import (
	"context"

	"github.com/google/uuid"
	tfe "github.com/hashicorp/go-tfe"
)

type OrganizationService interface {
	List(ctx context.Context, query string) ([]*tfe.Organization, error)
	Create(ctx context.Context, org *tfe.Organization) (*tfe.Organization, error)
	Read(ctx context.Context, name string) (*tfe.Organization, error)
	Update(ctx context.Context, name string, org *tfe.Organization) error
	Delete(ctx context.Context, name string) error
	GetIDByName(ctx context.Context, name string) (uuid.UUID, error)
}

type ProjectService interface {
	List(ctx context.Context, orgID uuid.UUID) ([]*tfe.Project, error)
	Create(ctx context.Context, project *tfe.Project) (*tfe.Project, error)
	Read(ctx context.Context, projectID string) (*tfe.Project, error)
	Update(ctx context.Context, project *tfe.Project) (*tfe.Project, error)
	Delete(ctx context.Context, projectID string) error
}

type UserService interface {
	List(ctx context.Context) ([]*tfe.User, error)
	Create(ctx context.Context, user *tfe.User) (*tfe.User, error)
	Read(ctx context.Context, userID string) (*tfe.User, error)
	Update(ctx context.Context, userID string, user *tfe.User) (*tfe.User, error)
	Delete(ctx context.Context, userID string) error
}
