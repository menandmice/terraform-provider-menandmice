data "menandmice_range" "rang" {
  name = "0.0.0.0/0"
}

output "range" {
  value = data.menandmice_range.rang
}
