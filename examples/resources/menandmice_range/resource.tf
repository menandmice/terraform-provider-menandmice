terraform {
  required_providers {
    menandmice = {
      # uncomment for terraform 0.13 and higher
      version = "~> 0.2.1",
      source  = "local/menandmice",
    }
  }
}
resource "menandmice_range" "range1" {
  # cidr = "192.168.2.0/24" # TODO test
  from        = "192.168.2.0"
  to          = "192.168.2.255"
  title       = "Test Terraform network"
  description = "Test"
  auto_assign = true
  locked      = true
}

