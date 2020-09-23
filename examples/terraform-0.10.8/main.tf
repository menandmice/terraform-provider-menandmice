terraform {
  required_providers {
    menandmice = {
      versions = ["0.2"],
    }
  }
}

provider "menandmice" {
  endpoint = "mandm.example.net"
  username = "rens"
  tls_verify= false
}


data "menandmice_dns_zone" "zone1" {
  name = "rens.nl."
  authority = "mandm.example.net."
}

resource menandmice_dns_zone zone2{
  name    = "test"
  dnsviewref = "DNSView/1"

}


output "zone2"{
  value = "${menandmice_dns_zone.zone2.name}"
}
