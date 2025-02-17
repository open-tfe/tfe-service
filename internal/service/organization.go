package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/go-tfe"
	"github.com/open-tfe/tfe-service/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type organizationService struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewOrganizationService(db *gorm.DB, logger *zap.Logger) OrganizationService {
	return &organizationService{db: db, logger: logger}
}

func (s *organizationService) List(ctx context.Context, query string) ([]*tfe.Organization, error) {
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

func (s *organizationService) Create(ctx context.Context, tfeOrg *tfe.Organization) (*tfe.Organization, error) {
	org := models.FromTFEOrganization(tfeOrg)
	s.logger.Debug("converting from TFE organization", zap.Any("organization", tfeOrg))

	if err := s.db.Create(org).Error; err != nil {
		s.logger.Error("failed to create organization", zap.Error(err))
		return nil, err
	}

	return org.ToTFE(), nil
}

func (s *organizationService) Read(ctx context.Context, name string) (*tfe.Organization, error) {
	var org models.Organization
	if err := s.db.Where("name = ?", name).First(&org).Error; err != nil {
		s.logger.Error("failed to read organization", zap.Error(err))
		return nil, err
	}

	s.logger.Debug("converting to TFE organization", zap.Any("organization", org))
	return org.ToTFE(), nil
}

func (s *organizationService) Update(ctx context.Context, name string, tfeOrg *tfe.Organization) error {
	org := models.FromTFEOrganization(tfeOrg)
	s.logger.Debug("converting from TFE organization", zap.Any("organization", tfeOrg))

	if err := s.db.Where("name = ?", name).Updates(org).Error; err != nil {
		s.logger.Error("failed to update organization", zap.Error(err))
		return err
	}
	return nil
}

func (s *organizationService) Delete(ctx context.Context, name string) error {
	if err := s.db.Where("name = ?", name).Delete(&models.Organization{}).Error; err != nil {
		s.logger.Error("failed to delete organization", zap.Error(err))
		return err
	}
	return nil
}

func (s *organizationService) GetIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	var org models.Organization
	if err := s.db.Where("name = ?", name).First(&org).Error; err != nil {
		s.logger.Error("failed to get organization ID", zap.Error(err))
		return uuid.Nil, err
	}
	return org.ID, nil
}
