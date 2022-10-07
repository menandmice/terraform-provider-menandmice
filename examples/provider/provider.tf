terraform {
  required_providers {
    menandmice = {
      source = "menandmice/menandmice",
    }
  }
}
provider "menandmice" {
  endpoint   = "https://micetro.example.net" # can also be set with MENANDMICE_ENDPOINT environment variable
  username   = "apiuser"                     # can also be set with MENANDMICE_USERNAME environment variable
  password   = "secret"                      # can also be set with MENANDMICE_PASSWORD environment variable
  tls_verify = false                         # can also be set with MENANDMICE_TLS_VERIFY environment variable
}
