package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-tfe"
	"gorm.io/gorm"
)

// TwoFactor represents the two-factor authentication settings for a user
type TwoFactor struct {
	Enabled  bool `gorm:"default:false" jsonapi:"attr,enabled"`
	Verified bool `gorm:"default:false" jsonapi:"attr,verified"`
}

// UserPermissions represents the permissions associated with a user
type UserPermissions struct {
	CanCreateOrganizations bool `gorm:"default:false" jsonapi:"attr,can-create-organizations"`
	CanChangeEmail         bool `gorm:"default:false" jsonapi:"attr,can-change-email"`
	CanChangeUsername      bool `gorm:"default:false" jsonapi:"attr,can-change-username"`
	CanManageUserTokens    bool `gorm:"default:false" jsonapi:"attr,can-manage-user-tokens"`
	CanView2FaSettings     bool `gorm:"default:false" jsonapi:"attr,can-view-2fa-settings"`
	CanManageHcpAccount    bool `gorm:"default:false" jsonapi:"attr,can-manage-hcp-account"`
}

type User struct {
	gorm.Model
	ID               uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" jsonapi:"primary,users"`
	AvatarURL        string           `jsonapi:"attr,avatar-url"`
	Email            string           `gorm:"uniqueIndex;not null" jsonapi:"attr,email"`
	IsServiceAccount bool             `gorm:"default:false" jsonapi:"attr,is-service-account"`
	TwoFactor        *TwoFactor       `gorm:"embedded;embeddedPrefix:two_factor_" jsonapi:"attr,two-factor"`
	UnconfirmedEmail string           `jsonapi:"attr,unconfirmed-email"`
	Username         string           `gorm:"uniqueIndex;not null" jsonapi:"attr,username"`
	IsSiteAdmin      bool             `gorm:"default:false" jsonapi:"attr,is-site-admin"`
	IsAdmin          bool             `gorm:"default:false" jsonapi:"attr,is-admin"`
	IsSsoLogin       bool             `gorm:"default:false" jsonapi:"attr,is-sso-login"`
	Permissions      *UserPermissions `gorm:"embedded;embeddedPrefix:permissions_" jsonapi:"attr,permissions"`
	LastLoginAt      time.Time        `jsonapi:"attr,last-login-at,iso8601"`
}

// ToTFE converts the internal User model to TFE format
func (u *User) ToTFE() *tfe.User {
	return &tfe.User{
		ID:               u.ID.String(),
		AvatarURL:        u.AvatarURL,
		Email:            u.Email,
		IsServiceAccount: u.IsServiceAccount,
		TwoFactor: &tfe.TwoFactor{
			Enabled:  u.TwoFactor.Enabled,
			Verified: u.TwoFactor.Verified,
		},
		UnconfirmedEmail: u.UnconfirmedEmail,
		Username:         u.Username,
		IsSiteAdmin:      &u.IsSiteAdmin,
		IsAdmin:          &u.IsAdmin,
		IsSsoLogin:       &u.IsSsoLogin,
		Permissions: &tfe.UserPermissions{
			CanCreateOrganizations: u.Permissions.CanCreateOrganizations,
			CanChangeEmail:         u.Permissions.CanChangeEmail,
			CanChangeUsername:      u.Permissions.CanChangeUsername,
			CanManageUserTokens:    u.Permissions.CanManageUserTokens,
			CanView2FaSettings:     u.Permissions.CanView2FaSettings,
			CanManageHcpAccount:    u.Permissions.CanManageHcpAccount,
		},
	}
}

// FromTFEUser converts a TFE User to internal model
func FromTFEUser(tfeUser *tfe.User) *User {
	id, _ := uuid.Parse(tfeUser.ID)
	return &User{
		ID:               id,
		AvatarURL:        tfeUser.AvatarURL,
		Email:            tfeUser.Email,
		IsServiceAccount: tfeUser.IsServiceAccount,
		TwoFactor: &TwoFactor{
			Enabled:  tfeUser.TwoFactor.Enabled,
			Verified: tfeUser.TwoFactor.Verified,
		},
		UnconfirmedEmail: tfeUser.UnconfirmedEmail,
		Username:         tfeUser.Username,
		IsAdmin:          *tfeUser.IsAdmin,
		IsSsoLogin:       *tfeUser.IsSsoLogin,
		Permissions: &UserPermissions{
			CanCreateOrganizations: tfeUser.Permissions.CanCreateOrganizations,
			CanChangeEmail:         tfeUser.Permissions.CanChangeEmail,
			CanChangeUsername:      tfeUser.Permissions.CanChangeUsername,
			CanManageUserTokens:    tfeUser.Permissions.CanManageUserTokens,
			CanView2FaSettings:     tfeUser.Permissions.CanView2FaSettings,
			CanManageHcpAccount:    tfeUser.Permissions.CanManageHcpAccount,
		},
	}
}
