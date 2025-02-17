table "projects" {
  schema = schema.public
  column "id" {
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "is_unified" {
    type = boolean
    default = false
  }
  column "name" {
    type = varchar(255)
    null = false
  }
  column "description" {
    type = text
    null = true
  }
  column "organization_id" {
    type = uuid
    null = false
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

  foreign_key "fk_projects_organization" {
    columns = [column.organization_id]
    ref_columns = [table.organizations.column.id]
    on_delete = CASCADE
  }

  index "idx_projects_name_org" {
    columns = [column.name, column.organization_id]
    unique = true
  }

  index "idx_projects_deleted_at" {
    columns = [column.deleted_at]
  }
}
