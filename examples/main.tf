terraform {
  required_providers {
    menandmice = {
      # uncomment for terraform 0.13 and higher
      # version = "~> 0.2",
      # source  = "local/menandmice",
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
  server = "mandm.example.net."
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
  name = "test"
  zone = "rens.nl."
  server = "mandm.example.net."
  type = "A"
}

resource menandmice_dns_record rec2 {
  name    = "test"
  zone    = "example.net."
  server  = "mandm.example.net."
  data    = "127.0.0.7"
  type    = "A"
}

data menandmice_ipam_record ipam1 {
  address = "2001:db8:0:0:0:0:0:25"
}

resource menandmice_ipam_record ipam2 {
  address = "2001:db8:0:0:0:0:0:29"
  custom_properties = {"location":"here"}
  claimed = true
}

resource menandmice_ipam_record ipam3 {
  free_ip {
    range = "172.16.17.0/24"
  }
  claimed = true
}

# // TODO can i make it so that this is possible
# locals {
#   free_ip_range ={
#     range = "172.16.17.0/24"
#     start_at = "172.16.17.100"
#     ping = true
#
#   }
# }
#
# resource menandmice_ipam_record ipam4 {
#   free_ip = local.free_ip_range
# }

data menandmice_dhcp_reservation reservation1 {
   name = "test"
}

resource menandmice_dhcp_reservation reservation2 {
  owner = "mandm.example.net."
  name    = "test5"
  client_identifier = "44:55:66:77:88:00"
  servername = "testname"
  next_server = "server1"
  reservation_method = "ClientIdentifier"
  # description = "test description"
  addresses = ["172.16.17.5","172.16.17.6"]
  ddns_hostname = "test"
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

output ipam1 {
  value = data.menandmice_ipam_record.ipam1
}

output ipam2 {
  value = menandmice_ipam_record.ipam2
}

output ipam3 {
  value = menandmice_ipam_record.ipam3
}

# output ipam4 {
#   value = menandmice_ipam_record.ipam4
# }

output reservation1 {
  value = data.menandmice_dhcp_reservation.reservation1
}
output reservation2 {
  value = menandmice_dhcp_reservation.reservation2
}
