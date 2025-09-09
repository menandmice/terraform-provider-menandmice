terraform {
  required_providers {
    menandmice = {
      source = "menandmice/menandmice",
    }
  }
}

data "menandmice_dhcp_reservation" "reservation1" {
  name = "reserved1"
}
