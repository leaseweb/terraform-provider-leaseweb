# An access key can be imported by specifying the object storage ID, the user ID and the access key ID in comma-separated format.
# Note that access_key & secret_access_key cannot be retrieved after creation, so both are null on an imported access key.
terraform import leaseweb_object_storage_access_key.example "5c9d2a1e-3b47-4f8d-9e21-7a6c8b0d4f33,4e7a9d02-6f31-48b5-a0c9-1d2e8f4b7a56,9b0c3e58-7a24-4d61-8f39-2e5a1c6b0d87"
