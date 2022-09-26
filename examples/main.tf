terraform {
  required_providers {
    menandmice = {
      source = "menandmice/menandmice",
    }
  }
}

provider "menandmice" {
  # endpoint = "https://micetro.example.net" # can also be set with MENANDMICE_ENDPOINT environment variable
  # username = "apiuser"           # can also be set with MENANDMICE_USERNAME environment variable
  # password = "secret"           # can also be set with MENANDMICE_PASSWORD environment variable
  tls_verify = false # can also be set with MENANDMICE_TLS_VERIFY environment variable
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

module "data-source_dhcp_scope" {
  source = "./data-sources/menandmice_dhcp_scope"
}
module "data-source_dns_record" {
  source = "./data-sources/menandmice_dns_record"
}

module "data-source_dns_zone" {
  source = "./data-sources/menandmice_dns_zone"
}
