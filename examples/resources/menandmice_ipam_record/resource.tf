terraform {
  required_providers {
    menandmice = {
      # uncomment for terraform 0.13 and higher
      version = "~> 0.2",
      source  = "local/menandmice",
    }
  }
}
resource menandmice_ipam_record ipam2 {
  address = "192.168.2.3"
  custom_properties = {"location":"here"}
  claimed = true
}
