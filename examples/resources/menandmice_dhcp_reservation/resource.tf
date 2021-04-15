terraform {
  required_providers {
    menandmice = {
      # uncomment for terraform 0.13 and higher
      version = "~> 0.2",
      source  = "local/menandmice",
    }
  }
}
resource menandmice_dhcp_reservation reservation1 {
  owner = "mandm.example.net."
  name    = "test1"
  client_identifier = "44:55:66:77:88:01"
  servername = "testname"
  next_server = "server1"
  reservation_method = "ClientIdentifier"
  # description = "test description" # only valid for some dhcp servers
  addresses = ["192.168.2.10"]
  ddns_hostname = "test"
}

data menandmice_dhcp_scope scope2{
  dhcp_server= "mandm.example.net."
  cidr = "192.168.2.0/24"
}

resource menandmice_dhcp_reservation reservation2 {
  owner = "mandm.example.net."
  name    = "test3"
  client_identifier = "44:55:66:77:88:01"
  servername = "testname"
  next_server = "server1"
  reservation_method = "ClientIdentifier"
  # description = "test description" # only valid for some dhcp servers
  addresses = ["192.168.2.11"]
  ddns_hostname = "test"
}
