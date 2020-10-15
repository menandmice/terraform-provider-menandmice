terraform {
  required_providers {
    menandmice = {
      versions = ["0.2"],
      #source  = "menandmice.com/menandmice/0.2",
    }
  }
}

provider menandmice {
  endpoint = "mandm.example.net"
  username = "rens"
  tls_verify= false
}


data menandmice_dns_zone zone1 {
  name = "rens.nl."
  authority = "mandm.example.net."
}

resource menandmice_dns_zone zone2{
  name    = "zone2."
  authority   = "mandm.example.net."
  adintegrated = false
  custom_properties = {"place" = "city","owner" = "me"}

  view = ""             # default ""
  type = "Master"       # default "Master"
  dnssecsigned = false  # default false
}

data menandmice_dns_record rec1 {
  fqdn = "test2.rens.nl."
}

resource menandmice_dns_record rec2 {
  name    = "test"
  zone    = "example.net."
  server  = "mandm.example.net."
  data    = "127.0.0.7"
  type    = "A"
}

output zone1{
  value = data.menandmice_dns_zone.zone1
}

output zone2{
  value = menandmice_dns_zone.zone2
}

output rec1{
  value = data.menandmice_dns_record.rec1
}

output rec2 {
  value = menandmice_dns_record.rec2
}
