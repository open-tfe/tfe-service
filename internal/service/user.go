package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/open-tfe/tfe-service/internal/constants"
	"github.com/open-tfe/tfe-service/internal/models"
	"go.uber.org/zap"
)

func (s *service) ListUsers(ctx context.Context) ([]*tfe.User, error) {
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

func (s *service) CreateUser(ctx context.Context, user *tfe.User) (*tfe.User, error) {
	dbUser := models.FromTFEUser(user)
	s.logger.Debug("converting from TFE user", zap.Any("user", user))

	if err := s.db.Create(dbUser).Error; err != nil {
		s.logger.Error("failed to create user", zap.Error(err))
		return nil, err
	}
	return dbUser.ToTFE(), nil
}

func (s *service) ReadUser(ctx context.Context, userID string) (*tfe.User, error) {
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

func (s *service) UpdateUser(ctx context.Context, userID string, user *tfe.User) (*tfe.User, error) {
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

func (s *service) DeleteUser(ctx context.Context, userID string) error {
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

func (s *service) ReadCurrentUser(ctx context.Context) (*tfe.User, error) {
	// Get the email from the context
	email, ok := ctx.Value(constants.UserEmailKey).(string)
	if !ok || email == "" {
		s.logger.Error("user email not found in context")
		return nil, fmt.Errorf("user email not found in context")
	}

	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		s.logger.Error("failed to read current user", zap.Error(err), zap.String("email", email))
		return nil, err
	}
	user.Permissions = &models.UserPermissions{
		CanCreateOrganizations: true,
		CanChangeEmail:         true,
		CanChangeUsername:      true,
		CanManageUserTokens:    true,
		CanView2FaSettings:     true,
		CanManageHcpAccount:    true,
	}
	user.TwoFactor = &models.TwoFactor{
		Enabled:  true,
		Verified: true,
	}
	s.logger.Debug("converting current user to TFE user", zap.Any("user", user))
	return user.ToTFE(), nil
}
