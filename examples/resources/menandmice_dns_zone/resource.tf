terraform {
  required_providers {
    menandmice = {
      # uncomment for terraform 0.13 and higher
      version = "~> 0.2",
      source  = "local/menandmice",
    }
  }
}
resource menandmice_dns_zone zone2{
  name    = "zone2.net."
  authority   = "micetro.example.net."
  adintegrated = false
  custom_properties = {"place" = "city","owner" = "me"}

  view = ""             # default ""
  type = "Master"       # default "Master"
  dnssecsigned = false  # default false
}
