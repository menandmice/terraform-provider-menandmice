terraform {
  required_providers {
    menandmice = {
      source = "menandmice/menandmice",
    }
  }
}
data "menandmice_dns_zone" "zone1" {
  name   = "zone1.net."
  server = "micetro.example.net."
}
