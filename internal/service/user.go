package service

import (
	"context"

	"github.com/google/uuid"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/open-tfe/tfe-service/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userService struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserService(db *gorm.DB, logger *zap.Logger) UserService {
	return &userService{db: db, logger: logger}
}

func (s *userService) List(ctx context.Context) ([]*tfe.User, error) {
	var users []*models.User
	if err := s.db.Find(&users).Error; err != nil {
		s.logger.Error("failed to list users", zap.Error(err))
		return nil, err
	}

	tfeUsers := make([]*tfe.User, len(users))
	for i, user := range users {
		s.logger.Debug("converting to TFE user", zap.Any("user", user))
		tfeUsers[i] = user.ToTFE()
	}
	return tfeUsers, nil
}

func (s *userService) Create(ctx context.Context, user *tfe.User) (*tfe.User, error) {
	dbUser := models.FromTFEUser(user)
	s.logger.Debug("converting from TFE user", zap.Any("user", user))

	if err := s.db.Create(dbUser).Error; err != nil {
		s.logger.Error("failed to create user", zap.Error(err))
		return nil, err
	}
	return dbUser.ToTFE(), nil
}

func (s *userService) Read(ctx context.Context, userID string) (*tfe.User, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := s.db.Where("id = ?", id).First(&user).Error; err != nil {
		s.logger.Error("failed to read user", zap.Error(err))
		return nil, err
	}
	s.logger.Debug("converting to TFE user", zap.Any("user", user))
	return user.ToTFE(), nil
}

func (s *userService) Update(ctx context.Context, userID string, user *tfe.User) (*tfe.User, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	dbUser := models.FromTFEUser(user)
	dbUser.ID = id
	s.logger.Debug("converting from TFE user", zap.Any("user", user))

	if err := s.db.Save(dbUser).Error; err != nil {
		s.logger.Error("failed to update user", zap.Error(err))
		return nil, err
	}
	return dbUser.ToTFE(), nil
}

func (s *userService) Delete(ctx context.Context, userID string) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	if err := s.db.Where("id = ?", id).Delete(&models.User{}).Error; err != nil {
		s.logger.Error("failed to delete user", zap.Error(err))
		return err
	}
	return nil
}
