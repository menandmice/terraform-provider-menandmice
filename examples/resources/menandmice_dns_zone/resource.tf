resource "menandmice_dns_zone" "zone2" {
  name              = "zone2.net."
  authority         = "micetro.example.net."
  custom_properties = { "place" = "city", "owner" = "me" }

  view          = ""       # default ""
  type          = "Master" # default "Master"
  dnssec_signed = false    # default false
}
