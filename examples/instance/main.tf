terraform {
  required_version = ">= 1.1.0"
  required_providers {
    elice = {
      source = "github.com/elice-dev/elice"
    }
  }
}

provider "elice" {
  host         = var.host
  token        = var.token
  organization = var.organization
}

resource "elice_instance" "example" {
  title            = "ex-instance"
  image_id         = "65f9e2f7-3c1a-4246-8e84-b70669d4c2fe"
  instance_type_id = "e4802a59-cca3-4cae-a719-5b3c8442c5d8"
  disk             = 128
}