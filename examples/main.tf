terraform {
  required_providers {
    menandmice = {
      versions = ["0.2"],
      #source  = "menandmice.com/menandmice/0.2",
    }
  }
}

provider menandmice {
  endpoint = "mandm.example.net"
  username = "rens"
}

data "menandmice_dnsrecord" "test2" {
  domain = "test2.rens.nl."
}

resource menandmice_dnsrecord rec1 {
  name    = "test"
  data    = "127.0.0.7"
  type    = "A"
  dnszone = "rens.nl"
}

output test1{

  value = data.menandmice_dnsrecord.test2

}

output rec1 {

  value = menandmice_dnsrecord.rec1
}
