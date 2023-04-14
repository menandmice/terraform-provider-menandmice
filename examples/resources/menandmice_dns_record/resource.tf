data "menandmice_dns_zone" "zone1" {
  name   = "zone1.net."
  server = "micetro.example.net."
}

resource "menandmice_dns_record" "rec1" {
  name   = "test1"
  zone   = data.menandmice_dns_zone.zone1.name   # "zone1.net."
  server = data.menandmice_dns_zone.zone1.server # "micetro.example.net."
  data   = "192.168.2.2"                         # this will asign/claim  "192.168.2.2" ipam records
  type   = "A"
}

resource "menandmice_dns_record" "rec2" {
  name         = "test2"
  dns_zone_ref = data.menandmice_dns_zone.zone1.ref
  data         = "192.168.2.2" # this will asign/claim  "192.168.2.2" ipam records
  type         = "A"
}
