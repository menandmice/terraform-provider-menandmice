resource "menandmice_ipam_record" "example1" {
  address           = "192.168.2.40"
  custom_properties = { "location" : "here" }
  claimed           = true
}

data "menandmice_range" "range1" {
  name = "192.168.2.0/24"
}


resource "menandmice_ipam_record" "example2" {
  free_ip {
    range    = data.menandmice_range.range1.name
    start_at = "192.168.2.50"
  }
}
