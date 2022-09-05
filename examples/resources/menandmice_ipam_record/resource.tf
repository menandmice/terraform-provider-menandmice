terraform {
  required_providers {
    menandmice = {
      # uncomment for terraform 0.13 and higher
      version = "~> 0.2",
      source  = "local/menandmice",
    }
  }
}
resource "menandmice_ipam_record" "ipam1" {
  address           = "192.168.2.3"
  custom_properties = { "location" : "here" }
  claimed           = true
}

resource "menandmice_ipam_record" "ipam2" {
  free_ip {
    range    = "192.168.2.0/24"
    start_at = "192.168.2.50"
  }
}
