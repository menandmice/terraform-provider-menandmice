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
  # endpoint = "mandm.example.net" # can also be set with MENANDMICE_ENDPOINT environment variable
  # username = "apiuser"           # can also be set with MENANDMICE_USERNAME environment variable
  # password = "secret"           # can also be set with MENANDMICE_PASSWORD environment variable
  tls_verify= false              # can also be set with MENANDMICE_TLS_VERIFY environment variable
}

module "resource_ipam_record" {
  source = "./resources/menandmice_ipam_record"
}

module "resource_dhcp_reservation" {
  source = "./resources/menandmice_dhcp_reservation"
}

module "resource_dns_record" {
  source = "./resources/menandmice_dns_record"
}

module "resource_dns_zone" {
  source = "./resources/menandmice_dns_zone"
}

module "data-source_ipam_record" {
  source = "./data-sources/menandmice_ipam_record"
}

module "data-source_dhcp_reservation" {
  source = "./data-sources/menandmice_dhcp_reservation"
}

module "data-source_dns_record" {
  source = "./data-sources/menandmice_dns_record"
}

module "data-source_dns_zone" {
  source = "./data-sources/menandmice_dns_zone"
}
