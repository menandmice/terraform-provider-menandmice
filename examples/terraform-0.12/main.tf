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
  name    = "test"
  authority   = "mandm.example.net."
  adintegrated = false
  type = "Master"
  # masters = ["::1"]
  # adreplicationtype = "None"
  dnssecsigned = true
}

data menandmice_dns_record rec1 {
  fqdn = "test2.rens.nl."
}

resource menandmice_dns_record rec2 {
  name    = "test"
  data    = "127.0.0.7"
  type    = "A"
  dns_zone_ref =  data.menandmice_dns_zone.zone1.ref
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
