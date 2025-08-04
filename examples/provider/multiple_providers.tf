terraform {
  required_providers {
    leaseweb = {
      version = "~> 1.2.0"
      source  = "leaseweb/leaseweb"
    }
  }
}

provider "leaseweb" {
  alias = "nl"
  token = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
}

provider "leaseweb" {
  alias = "us"
  token = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
}

resource "leaseweb_dedicated_server" "web-nl" {
  provider  = leaseweb.nl
  reference = "web-nl"
}

resource "leaseweb_dedicated_server" "web-us" {
  provider  = leaseweb.us
  reference = "web-us"
}
