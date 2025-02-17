table "users" {
  schema = schema.public
  column "id" {
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "username" {
    type = varchar(255)
    null = false
  }
  column "email" {
    type = varchar(255)
    null = false
  }
  column "password" {
    type = varchar(255)
    null = false
  }
  column "is_service_account" {
    type = boolean
    default = false
  }
  column "is_admin" {
    type = boolean
    default = false
  }
  column "is_site_admin" {
    type = boolean
    default = false
  }
  column "two_factor_enabled" {
    type = boolean
    default = false
  }
  column "two_factor_verified" {
    type = boolean
    default = false
  }
  column "avatar_url" {
    type = text
    null = true
  }
  column "enterprise_support" {
    type = boolean
    default = false
  }
  column "last_login_at" {
    type = timestamp
    null = true
  }
  column "unconfirmed_email" {
    type = varchar(255)
    null = true
  }
  column "has_git_hub_app_token" {
    type = boolean
    default = false
  }
  column "onboarding_status" {
    type = varchar(255)
    null = true
  }
  column "created_at" {
    type = timestamp
    default = sql("NOW()")
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

  index "idx_users_username" {
    columns = [column.username]
    unique = true
  }

  index "idx_users_email" {
    columns = [column.email]
    unique = true
  }

  index "idx_users_deleted_at" {
    columns = [column.deleted_at]
  }
} 