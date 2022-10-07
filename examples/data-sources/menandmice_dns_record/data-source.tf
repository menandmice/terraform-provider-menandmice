data "menandmice_dns_zone" "zone1" {
  name   = "zone1.net."
  server = "micetro.example.net."
}

data "menandmice_dns_record" "rec1" {
  name   = "test"
  zone   = data.menandmice_dns_zone.zone1.name # "zone1.net."
  server = "micetro.example.net."
  type   = "A"
}
