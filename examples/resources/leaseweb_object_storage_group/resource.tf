# Manage an object storage group
resource "leaseweb_object_storage_group" "readonly" {
  object_storage_id = "5c9d2a1e-3b47-4f8d-9e21-7a6c8b0d4f33"
  display_name      = "Read only"
  unique_name       = "read-only"

  s3_policies = jsonencode({
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["s3:GetObject", "s3:ListBucket"]
        Resource = ["arn:aws:s3:::*"]
      }
    ]
  })
}
