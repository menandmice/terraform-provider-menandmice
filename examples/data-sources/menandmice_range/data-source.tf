terraform {
  required_providers {
    menandmice = {
      # uncomment for terraform 0.13 and higher
      version = "~> 0.2.1",
      source  = "local/menandmice",
    }
  }
}

data "menandmice_range" "rang" {
  name = "0.0.0.0/0"
}

output "range" {
  value = data.menandmice_range.rang
}
