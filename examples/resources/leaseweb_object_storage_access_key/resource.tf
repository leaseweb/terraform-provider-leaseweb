# Manage an object storage access key
resource "leaseweb_object_storage_access_key" "backup_agent" {
  object_storage_id = "5c9d2a1e-3b47-4f8d-9e21-7a6c8b0d4f33"
  user_id           = leaseweb_object_storage_user.backup_agent.id
}

# An access key that expires
resource "leaseweb_object_storage_access_key" "temporary" {
  object_storage_id = "5c9d2a1e-3b47-4f8d-9e21-7a6c8b0d4f33"
  user_id           = leaseweb_object_storage_user.backup_agent.id
  expires_at        = "2027-01-01T00:00:00Z"
}

# access_key & secret_access_key are only returned when the key is created,
# so store them somewhere safe on the first apply.
output "backup_agent_secret_access_key" {
  value     = leaseweb_object_storage_access_key.backup_agent.secret_access_key
  sensitive = true
}
