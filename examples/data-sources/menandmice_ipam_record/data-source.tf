terraform {
  required_providers {
    menandmice = {
      # uncomment for terraform 0.13 and higher
      version = "~> 0.2",
      source  = "local/menandmice",
    }
  }
}
data "menandmice_ipam_record" "ipam1" {
  address = "192.168.2.2"
}
