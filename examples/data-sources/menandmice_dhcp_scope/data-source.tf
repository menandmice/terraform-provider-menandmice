data "menandmice_dhcp_scope" "scope1" {
  dhcp_server = "micetro.example.net."
  cidr        = "192.168.2.0/24"
}
