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
