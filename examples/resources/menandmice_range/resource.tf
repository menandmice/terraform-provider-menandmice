resource "menandmice_range" "example1" {
  cidr   = "192.168.5.0/24"
  title  = "Test Terraform example1"
  subnet = true
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
    mask  = 24
  }
  title       = "Test Terraform example3"
  description = "Test"
}

resource "menandmice_range" "example4" {
  free_range {
    ranges = ["10.0.16.0/24","10.0.17.0/24"]
    size = 100
  }
  title       = "Test Terraform example3"
  description = "Test"
}
