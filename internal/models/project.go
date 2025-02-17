package models

import (
	"github.com/google/uuid"
	"github.com/hashicorp/go-tfe"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	ID             string        `gorm:"primaryKey;type:varchar(255)" jsonapi:"primary,projects"`
	IsUnified      bool          `gorm:"default:false" jsonapi:"attr,is-unified"`
	Name           string        `gorm:"not null" jsonapi:"attr,name"`
	Description    string        `gorm:"type:text" jsonapi:"attr,description"`
	OrganizationID uuid.UUID     `gorm:"type:uuid;not null"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID" jsonapi:"relation,organization"`
	// EffectiveTagBindings []*EffectiveTagBinding `jsonapi:"relation,effective-tag-bindings"`
}

// ToTFE converts the internal Project model to TFE format
func (p *Project) ToTFE() *tfe.Project {
	return &tfe.Project{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		IsUnified:   p.IsUnified,
	}
}

// FromTFE converts a TFE Project to internal model
func FromTFEProject(proj *tfe.Project) *Project {
	return &Project{
		ID:          proj.ID,
		Name:        proj.Name,
		Description: proj.Description,
		IsUnified:   proj.IsUnified,
	}
}
