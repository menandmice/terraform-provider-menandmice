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


data "menandmice_dnszone" "zone1" {
  domain = "rens.nl."
}

# data "menandmice_dnsrecord" "rec1" {
#   fqdn = "test2.rens.nl."
# }
#
# resource menandmice_dnsrecord rec2 {
#   name    = "test"
#   data    = "127.0.0.7"
#   type    = "A"
#   dnszone = "rens.nl"
# }
#
output zone1{
  value = data.menandmice_dnszone.zone1
}

# output rec1{
#   value = data.menandmice_dnsrecord.rec1
# }
#
# output rec2 {
#   value = menandmice_dnsrecord.rec2
# }
