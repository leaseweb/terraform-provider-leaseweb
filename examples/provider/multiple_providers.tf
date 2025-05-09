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
  token = "527070ca-8449-4f06-b609-ec6797bd8222"
}

provider "leaseweb" {
  alias = "us"
  token = "416fa444-5e96-4198-a4f7-297cbbc3cc70"
}

resource "leaseweb_dedicated_server" "web-nl" {
  provider  = leaseweb.nl
  reference = "web-nl"
}

resource "leaseweb_dedicated_server" "web-us" {
  provider  = leaseweb.us
  reference = "web-us"
}