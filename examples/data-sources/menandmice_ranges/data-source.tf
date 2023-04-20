data "menandmice_ranges" "rangs" {
  # not requestion all ranges, but first 10 can improve speed if there exist a lot ranges
  limit        = 10
  folder       = "testfolder"
  subnet       = true
  is_container = false
}

output "ranges" {
  value = data.menandmice_range.rangs
}
