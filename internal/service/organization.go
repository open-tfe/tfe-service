package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/go-tfe"
	"github.com/open-tfe/tfe-service/internal/models"
	"go.uber.org/zap"
)

func (s *service) ListOrganizations(ctx context.Context, query string) ([]*tfe.Organization, error) {
	var orgs []*models.Organization
	db := s.db

	if query != "" {
		db = db.Where("name ILIKE ? OR email ILIKE ?", "%"+query+"%", "%"+query+"%")
	}

	if err := db.Find(&orgs).Error; err != nil {
		return nil, err
	}

	tfeOrgs := make([]*tfe.Organization, len(orgs))
	for i, org := range orgs {
		s.logger.Debug("converting to TFE organization", zap.Any("organization", org))
		tfeOrgs[i] = org.ToTFE()
	}

	return tfeOrgs, nil
}

func (s *service) CreateOrganization(ctx context.Context, org *tfe.Organization) (*tfe.Organization, error) {
	tforg := models.FromTFEOrganization(org)
	s.logger.Debug("converting from TFE organization", zap.Any("organization", org))

	if err := s.db.Create(tforg).Error; err != nil {
		s.logger.Error("failed to create organization", zap.Error(err))
		return nil, err
	}

	return tforg.ToTFE(), nil
}

func (s *service) ReadOrganization(ctx context.Context, name string) (*tfe.Organization, error) {
	var org models.Organization
	if err := s.db.Where("name = ?", name).First(&org).Error; err != nil {
		s.logger.Error("failed to read organization", zap.Error(err))
		return nil, err
	}
	projects, _, err := s.ListProjects(ctx, org.ID)
	if err != nil {
		s.logger.Error("failed to list projects", zap.Error(err))
		return nil, err
	}
	org.Projects = projects
	if org.DefaultProject == nil {
		org.DefaultProject = projects[0]
	}

	s.logger.Debug("converting to TFE organization", zap.Any("organization", org))
	return org.ToTFE(), nil
}

func (s *service) UpdateOrganization(ctx context.Context, name string, org *tfe.Organization) error {
	tforg := models.FromTFEOrganization(org)
	s.logger.Debug("converting from TFE organization", zap.Any("organization", org))

	if err := s.db.Where("name = ?", name).Updates(tforg).Error; err != nil {
		s.logger.Error("failed to update organization", zap.Error(err))
		return err
	}
	return nil
}

func (s *service) DeleteOrganization(ctx context.Context, name string) error {
	if err := s.db.Where("name = ?", name).Delete(&models.Organization{}).Error; err != nil {
		s.logger.Error("failed to delete organization", zap.Error(err))
		return err
	}
	return nil
}

func (s *service) GetOrganizationIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	var org models.Organization
	if err := s.db.Where("name = ?", name).First(&org).Error; err != nil {
		s.logger.Error("failed to get organization ID", zap.Error(err))
		return uuid.Nil, err
	}
	return org.ID, nil
}

func (s *service) ReadOrganizationEntitlements(ctx context.Context, name string) (*tfe.Entitlements, error) {

	entitlements := &tfe.Entitlements{
		Agents:                     true,
		AuditLogging:               true,
		CostEstimation:             true,
		GlobalRunTasks:             true,
		Operations:                 true,
		PrivateModuleRegistry:      true,
		RunTasks:                   true,
		SSO:                        true,
		Sentinel:                   true,
		StateStorage:               true,
		Teams:                      true,
		VCSIntegrations:            true,
		WaypointActions:            true,
		WaypointTemplatesAndAddons: true,
	}

	return entitlements, nil
}
