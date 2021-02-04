data menandmice_dns_zone zone1 {
  name = "zone1.net."
  server = "mandm.example.net."
}

data menandmice_dns_record rec1 {
  name = "test"
  zone = data.menandmice_dns_zone.zone1.name  # "zone1.net."
  server = "mandm.example.net."
  type = "A"
}
