package service

import (
	"context"

	"github.com/google/uuid"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/open-tfe/tfe-service/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type projectService struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewProjectService(db *gorm.DB, logger *zap.Logger) ProjectService {
	return &projectService{db: db, logger: logger}
}

func (s *projectService) List(ctx context.Context, orgID uuid.UUID) ([]*tfe.Project, error) {
	var projects []*models.Project
	if err := s.db.Where("organization_id = ?", orgID).Find(&projects).Error; err != nil {
		s.logger.Error("failed to list projects", zap.Error(err))
		return nil, err
	}

	// Convert to TFE format
	tfeProjects := make([]*tfe.Project, len(projects))
	for i, proj := range projects {
		s.logger.Debug("converting to TFE project", zap.Any("project", proj))
		tfeProjects[i] = proj.ToTFE()
	}
	return tfeProjects, nil
}

func (s *projectService) Create(ctx context.Context, project *tfe.Project) (*tfe.Project, error) {
	dbProject := models.FromTFEProject(project)
	s.logger.Debug("converting from TFE project", zap.Any("project", project))

	if err := s.db.Create(dbProject).Error; err != nil {
		s.logger.Error("failed to create project", zap.Error(err))
		return nil, err
	}
	return dbProject.ToTFE(), nil
}

func (s *projectService) Read(ctx context.Context, projectID string) (*tfe.Project, error) {
	var project models.Project
	if err := s.db.Where("id = ?", projectID).First(&project).Error; err != nil {
		s.logger.Error("failed to read project", zap.Error(err))
		return nil, err
	}
	s.logger.Debug("converting to TFE project", zap.Any("project", project))
	return project.ToTFE(), nil
}

func (s *projectService) Update(ctx context.Context, project *tfe.Project) (*tfe.Project, error) {
	dbProject := models.FromTFEProject(project)
	s.logger.Debug("converting from TFE project", zap.Any("project", project))

	if err := s.db.Save(dbProject).Error; err != nil {
		s.logger.Error("failed to update project", zap.Error(err))
		return nil, err
	}
	return dbProject.ToTFE(), nil
}

func (s *projectService) Delete(ctx context.Context, projectID string) error {
	if err := s.db.Where("id = ?", projectID).Delete(&models.Project{}).Error; err != nil {
		s.logger.Error("failed to delete project", zap.Error(err))
		return err
	}
	return nil
}
