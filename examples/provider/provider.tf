terraform {
  required_providers {
    leaseweb = {
      version = "~> 1.2.0"
      source  = "leaseweb/leaseweb"
    }
  }
}

provider "leaseweb" {
  token = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
}
