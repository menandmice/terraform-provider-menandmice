# import with dnszone ref
terraform import menandmice_dns_zone.resourcename DNSZones/659

# import with readable name
terraform import menandmice_dns_zone.resourcename mandm.example.net::zone1  #<server>:<view>:<dnzzone name>

