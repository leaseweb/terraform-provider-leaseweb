# Manage an object storage user
resource "leaseweb_object_storage_user" "backup_agent" {
  object_storage_id = "5c9d2a1e-3b47-4f8d-9e21-7a6c8b0d4f33"
  full_name         = "Backup agent"
  unique_name       = "backup-agent"
  groups            = [leaseweb_object_storage_group.readonly.id]
}
