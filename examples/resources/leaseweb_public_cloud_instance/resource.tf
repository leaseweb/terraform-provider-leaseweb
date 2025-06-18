# Manage example Public Cloud instance
resource "leaseweb_public_cloud_instance" "example" {
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
  image = {
    id = "UBUNTU_22_04_64BIT"
  }
  reference              = "my webserver"
  region                 = "eu-west-3"
  root_disk_storage_type = "CENTRAL"
  type                   = "lsw.m3.large"
  private_network        = true
}
