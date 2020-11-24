terraform {
  required_providers {
    menandmice = {
      # # uncomment for terraform 0.13 and higher
      # version = "~> 0.2",
      # source  = "local/menandmice",
    }
  }
}

provider menandmice {
  endpoint = "mandm.example.net"
  # username = "apiuser"
  tls_verify= false
}

data menandmice_dns_zone zone1 {
  name = "zone1.net."
  server = "mandm.example.net."
}

resource menandmice_dns_zone zone2{
  name    = "zone2.net."
  authority   = "mandm.example.net."
  adintegrated = false
  custom_properties = {"place" = "city","owner" = "me"}

  view = ""             # default ""
  type = "Master"       # default "Master"
  dnssecsigned = false  # default false
}

data menandmice_dns_record rec1 {
  name = "test"
  zone = data.menandmice_dns_zone.zone1.name  # "zone1.net."
  server = "mandm.example.net."
  type = "A"
}

resource menandmice_dns_record rec2 {
  name    = "test"
  zone    = menandmice_dns_zone.zone2.name      # "zone2.net."
  server  = "mandm.example.net."
  data    = "192.168.2.2" # this will asign/claim  "192.168.2.2" ipam records
  type    = "A"
}

data menandmice_ipam_record ipam1 {
  address = "192.168.2.2"
}

resource menandmice_ipam_record ipam2 {
  address = "192.168.2.3"
  custom_properties = {"location":"here"}
  claimed = true
}

resource menandmice_ipam_record ipam3 {
  free_ip {
    range = "192.168.2.0/24"
    start_at = "192.168.2.50"
  }
}

data menandmice_dhcp_reservation reservation1 {
   name = "reserved1"
}

resource menandmice_dhcp_reservation reservation2 {
  owner = "mandm.example.net."
  name    = "test5"
  client_identifier = "44:55:66:77:88:01"
  servername = "testname"
  next_server = "server1"
  reservation_method = "ClientIdentifier"
  # description = "test description" # only valid for some dhcp servers
  addresses = ["192.168.2.10"]
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

output reservation1 {
  value = data.menandmice_dhcp_reservation.reservation1
}

output reservation2 {
  value = menandmice_dhcp_reservation.reservation2
}
