terraform {
  required_providers {
    menandmice = {
      # uncomment for terraform 0.13 and higher
      version = "~> 0.2",
      source  = "local/menandmice",
    }
  }
}
data "menandmice_dhcp_reservation" "reservation1" {
  name = "reserved1"
}
