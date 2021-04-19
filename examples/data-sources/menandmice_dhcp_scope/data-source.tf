terraform {
  required_providers {
    menandmice = {
      # uncomment for terraform 0.13 and higher
      version = "~> 0.2",
      source  = "local/menandmice",
    }
  }
}

data menandmice_dhcp_scope scope1{
  dhcp_server= "micetro.example.net."
  cidr = "192.168.2.0/24"
}
