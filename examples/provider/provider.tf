terraform {
  required_providers {
    menandmice = {
      # # uncomment for terraform 0.13 and higher
      #version = "~> 0.2",
      source  = "local/menandmice",
    }
  }
}
provider menandmice {
  endpoint = "mandm.example.net" # can also be set with MENANDMICE_ENDPOINT environment variable
  username = "apiuser"           # can also be set with MENANDMICE_USERNAME environment variable
  password = "secret"           # can also be set with MENANDMICE_PASSWORD environment variable
  tls_verify= false              # can also be set with MENANDMICE_TLS_VERIFY environment variable
}
