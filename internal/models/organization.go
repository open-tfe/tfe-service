package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-tfe"
	"gorm.io/gorm"
)

type Organization struct {
	gorm.Model
	ID                  uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" jsonapi:"primary,organizations"`
	Name                string    `gorm:"uniqueIndex;not null" jsonapi:"attr,name"`
	AssessmentsEnforced bool      `gorm:"default:false" jsonapi:"attr,assessments-enforced"`
	// CollaboratorAuthPolicy                            AuthPolicyType           `jsonapi:"attr,collaborator-auth-policy"`
	CostEstimationEnabled bool      `gorm:"default:false" jsonapi:"attr,cost-estimation-enabled"`
	CreatedAt             time.Time `jsonapi:"attr,created-at,iso8601"`
	DefaultExecutionMode  string    `gorm:"type:varchar(255)" jsonapi:"attr,default-execution-mode"`
	Email                 string    `gorm:"not null" jsonapi:"attr,email"`
	ExternalID            string    `jsonapi:"attr,external-id"`
	IsUnified             bool      `gorm:"default:false" jsonapi:"attr,is-unified"`
	OwnersTeamSAMLRoleID  string    `jsonapi:"attr,owners-team-saml-role-id"`
	// Permissions                                       *OrganizationPermissions `jsonapi:"attr,permissions"`
	SAMLEnabled                                       bool      `gorm:"default:false" jsonapi:"attr,saml-enabled"`
	SessionRemember                                   int       `gorm:"default:20160" jsonapi:"attr,session-remember"`
	SessionTimeout                                    int       `gorm:"default:20160" jsonapi:"attr,session-timeout"`
	TrialExpiresAt                                    time.Time `jsonapi:"attr,trial-expires-at,iso8601"`
	TwoFactorConformant                               bool      `gorm:"default:false" jsonapi:"attr,two-factor-conformant"`
	SendPassingStatusesForUntriggeredSpeculativePlans bool      `gorm:"default:false" jsonapi:"attr,send-passing-statuses-for-untriggered-speculative-plans"`
	RemainingTestableCount                            int       `gorm:"default:0" jsonapi:"attr,remaining-testable-count"`
	SpeculativePlanManagementEnabled                  bool      `gorm:"default:false" jsonapi:"attr,speculative-plan-management-enabled"`
	// Optional: If enabled, SendPassingStatusesForUntriggeredSpeculativePlans needs to be false.
	AggregatedCommitStatusEnabled bool `gorm:"default:false" jsonapi:"attr,aggregated-commit-status-enabled,omitempty"`
	// Note: This will be false for TFE versions older than v202211, where the setting was introduced.
	// On those TFE versions, safe delete does not exist, so ALL deletes will be force deletes.
	AllowForceDeleteWorkspaces bool `gorm:"default:false" jsonapi:"attr,allow-force-delete-workspaces"`

	// Relations
	DefaultProject *Project  `gorm:"foreignKey:OrganizationID" jsonapi:"relation,default-project"`
	Projects       []Project `gorm:"foreignKey:OrganizationID"`
	// DefaultAgentPool *AgentPool `jsonapi:"relation,default-agent-pool"`

	// Deprecated: Use DataRetentionPolicyChoice instead.
	// DataRetentionPolicy *DataRetentionPolicy

	// **Note: This functionality is only available in Terraform Enterprise.**
	// DataRetentionPolicyChoice *DataRetentionPolicyChoice `jsonapi:"polyrelation,data-retention-policy"`
}

// ToTFE converts the internal Organization model to TFE format
func (o *Organization) ToTFE() *tfe.Organization {
	return &tfe.Organization{
		Name:                  o.Name,
		AssessmentsEnforced:   o.AssessmentsEnforced,
		CostEstimationEnabled: o.CostEstimationEnabled,
		CreatedAt:             o.CreatedAt,
		DefaultExecutionMode:  o.DefaultExecutionMode,
		Email:                 o.Email,
		ExternalID:            o.ExternalID,
		IsUnified:             o.IsUnified,
		OwnersTeamSAMLRoleID:  o.OwnersTeamSAMLRoleID,
		SAMLEnabled:           o.SAMLEnabled,
		SessionRemember:       o.SessionRemember,
		SessionTimeout:        o.SessionTimeout,
		TrialExpiresAt:        o.TrialExpiresAt,
		TwoFactorConformant:   o.TwoFactorConformant,
		SendPassingStatusesForUntriggeredSpeculativePlans: o.SendPassingStatusesForUntriggeredSpeculativePlans,
		RemainingTestableCount:                            o.RemainingTestableCount,
		SpeculativePlanManagementEnabled:                  o.SpeculativePlanManagementEnabled,
		AggregatedCommitStatusEnabled:                     o.AggregatedCommitStatusEnabled,
		AllowForceDeleteWorkspaces:                        o.AllowForceDeleteWorkspaces,
		DefaultProject:                                    o.DefaultProject.ToTFE(),
		// CollaboratorAuthPolicy: o.CollaboratorAuthPolicy,
		// Add other fields as needed
	}
}

// FromTFEOrganization converts a TFE Organization to internal model
func FromTFEOrganization(tfeOrg *tfe.Organization) *Organization {
	return &Organization{
		// ID:                    id,
		Name:                  tfeOrg.Name,
		AssessmentsEnforced:   tfeOrg.AssessmentsEnforced,
		CostEstimationEnabled: tfeOrg.CostEstimationEnabled,
		CreatedAt:             tfeOrg.CreatedAt,
		DefaultExecutionMode:  tfeOrg.DefaultExecutionMode,
		Email:                 tfeOrg.Email,
		ExternalID:            tfeOrg.ExternalID,
		IsUnified:             tfeOrg.IsUnified,
		OwnersTeamSAMLRoleID:  tfeOrg.OwnersTeamSAMLRoleID,
		SAMLEnabled:           tfeOrg.SAMLEnabled,
		SessionTimeout:        tfeOrg.SessionTimeout,
		SessionRemember:       tfeOrg.SessionRemember,
		TrialExpiresAt:        tfeOrg.TrialExpiresAt,
		TwoFactorConformant:   tfeOrg.TwoFactorConformant,
		SendPassingStatusesForUntriggeredSpeculativePlans: tfeOrg.SendPassingStatusesForUntriggeredSpeculativePlans,
		RemainingTestableCount:                            tfeOrg.RemainingTestableCount,
		SpeculativePlanManagementEnabled:                  tfeOrg.SpeculativePlanManagementEnabled,
		AggregatedCommitStatusEnabled:                     tfeOrg.AggregatedCommitStatusEnabled,
		AllowForceDeleteWorkspaces:                        tfeOrg.AllowForceDeleteWorkspaces,
		DefaultProject:                                    FromTFEProject(tfeOrg.DefaultProject),
		// CollaboratorAuthPolicy: tfeOrg.CollaboratorAuthPolicy,
		// Add other fields as needed
	}
}
