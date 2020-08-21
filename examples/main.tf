terraform {
  required_providers {
    menandmice = {
      versions = ["0.2"],
      #source  = "menandmice.com/menandmice/0.2",
    }
  }
}

provider menandmice {
      web      = "mandm.example.nett"
      username = "rens"
}

data "menandmice_dnsrecord" "tonk" {
  domain = "tonk.example.net"
}

resource menandmice_dnsrecord rec1 {
    name    = "test1"
    data    = "127.0.0.1"
    type    = "A"
    dnszone = "rens.nl"
}

output data_tonk {

  value = data.menandmice_dnsrecord.tonk
}

output rec1 {

  value = menandmice_dnsrecord.rec1
}
