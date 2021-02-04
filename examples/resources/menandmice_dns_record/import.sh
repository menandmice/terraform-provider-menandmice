# import with DNS record ref
terraform import menandmice_dns_record.resourcename DNSRecords/2294

# import with readable name
terraform import menandmice_dns_record.resourcename mandm.example.net.::www.test.org.:A # <dns server>:<view>:<zone>:<type>

