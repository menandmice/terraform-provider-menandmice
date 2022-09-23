terraform {
  required_providers {
    menandmice = {
      # uncomment for terraform 0.13 and higher
      version = "~> 0.2.1",
      source  = "local/menandmice",
    }
  }
}

resource "menandmice_range" "example1" {
  cidr  = "192.168.5.0/25"
  title = "Test Terraform example1"
}

resource "menandmice_range" "example2" {
  from        = "192.168.2.0"
  to          = "192.168.2.255"
  title       = "Test Terraform example2"
  description = "Test"
  auto_assign = true
  locked      = true
}

data "menandmice_range" "super_range" {
  name = "192.168.0.0/16"
}

resource "menandmice_range" "example3" {
  free_range {
    range = data.menandmice_range.super_range.name
    size =  255
  }
  title       = "Test Terraform example3"
  description = "Test"
}

