# Manage an object storage bucket
resource "leaseweb_object_storage_bucket" "backups" {
  object_storage_id     = "5c9d2a1e-3b47-4f8d-9e21-7a6c8b0d4f33"
  name                  = "my-backups"
  is_versioning_enabled = true
}

# A bucket with a quota of 500 GB
resource "leaseweb_object_storage_bucket" "archive" {
  object_storage_id     = "5c9d2a1e-3b47-4f8d-9e21-7a6c8b0d4f33"
  name                  = "my-archive"
  is_versioning_enabled = false
  quota                 = 500
}
