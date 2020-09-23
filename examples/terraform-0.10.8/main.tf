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


data "menandmice_dnszone" "zone1" {
  name = "rens.nl."
  authority = "mandm.example.net."
}

resource menandmice_dnszone zone2{
  name    = "test"
  dnsviewref = "DNSView/1"

}


output "zone2"{
  value = "${menandmice_dnszone.zone2.name}"
}
