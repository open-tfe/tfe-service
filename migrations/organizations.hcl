table "organizations" {
  schema = schema.public
  column "id" {
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "name" {
    type = varchar(255)
    null = false
  }
  column "assessments_enforced" {
    type = boolean
    default = false
  }
  column "cost_estimation_enabled" {
    type = boolean
    default = false
  }
  column "created_at" {
    type = timestamp
    default = sql("NOW()")
  }
  column "default_execution_mode" {
    type = varchar(255)
    null = true
  }
  column "email" {
    type = varchar(255)
    null = false
  }
  column "external_id" {
    type = varchar(255)
    null = true
  }
  column "is_unified" {
    type = boolean
    default = false
  }
  column "owners_team_saml_role_id" {
    type = varchar(255)
    null = true
  }
  column "saml_enabled" {
    type = boolean
    default = false
  }
  column "session_remember" {
    type = integer
    default = 20160 // 14 days in minutes
  }
  column "session_timeout" {
    type = integer
    default = 20160 // 14 days in minutes
  }
  column "trial_expires_at" {
    type = timestamp
    null = true
  }
  column "two_factor_conformant" {
    type = boolean
    default = false
  }
  column "send_passing_statuses_for_untriggered_speculative_plans" {
    type = boolean
    default = false
  }
  column "remaining_testable_count" {
    type = integer
    default = 0
  }
  column "speculative_plan_management_enabled" {
    type = boolean
    default = false
  }
  column "aggregated_commit_status_enabled" {
    type = boolean
    default = false
  }
  column "allow_force_delete_workspaces" {
    type = boolean
    default = false
  }
  column "updated_at" {
    type = timestamp
    default = sql("NOW()")
  }
  column "deleted_at" {
    type = timestamp
    null = true
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_organizations_name" {
    columns = [column.name]
    unique = true
  }

  index "idx_organizations_deleted_at" {
    columns = [column.deleted_at]
  }
}
