package service

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Services struct {
	Organization OrganizationService
	Project      ProjectService
	User         UserService
	logger       *zap.Logger
}

func NewServices(db *gorm.DB, logger *zap.Logger) *Services {
	return &Services{
		Organization: NewOrganizationService(db, logger),
		Project:      NewProjectService(db, logger),
		User:         NewUserService(db, logger),
		logger:       logger.With(zap.String("component", "services")),
	}
}
