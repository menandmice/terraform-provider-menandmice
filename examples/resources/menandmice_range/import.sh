# import with range ref
terraform import menandmice_range.resourcename Range/0

# import with readable name
terraform import menandmice_range.range1 192.168.2.0/24
# or
terraform import menandmice_range.range1 192.168.2.0-192.168.2.254
