terraform {
  required_providers {
    menandmice = {
      # uncomment for terraform 0.13 and higher
      version = "~> 0.2",
      source  = "local/menandmice",
    }
  }
}

data "menandmice_dns_zone" "zone1" {
  name   = "zone1.net."
  server = "micetro.example.net."
}

resource "menandmice_dns_record" "rec2" {
  name   = "test"
  zone   = data.menandmice_dns_zone.zone1.name   # "zone1.net."
  server = data.menandmice_dns_zone.zone1.server # "micetro.example.net."
  data   = "192.168.2.2"                         # this will asign/claim  "192.168.2.2" ipam records
  type   = "A"
}
