data "menandmice_dns_zones" "zones1" {
  limit  = 10
  folder = "AWS"
}

resource "menandmice_dns_record" "recs" {
  for_each = { for zone in data.menandmice_dns_zones.folder.zones1 : zone.name => zone.authority }
  name     = "test"
  zone     = each.key
  server   = each.value
  data     = "192.168.2.2" # this will asign/claim  "192.168.2.2" ipam records
  type     = "A"
}


data "menandmice_dns_zones" "zones2" {
  custom_properties = { "Owner" = "me" }
}

data "menandmice_dns_zones" "zones" {
  folder = "AD"
  type   = "Master"
  server = "micetro.example.net."
}



