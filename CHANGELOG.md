
# bet practice from https://www.terraform.io/plugin/sdkv2/best-practices/versioning

## 0.4.0 (Unreleased)

BREAKING CHANGES:

* data_source/dhcpscope cidr is renamed to range
* dnsrecord name attribute can't end with a "." anymore
* resouce/dnszone displayname became read only. setting it was allowed before, but did not work
* resouce/dnszone dynamic became read only. setting it was allowed before, but did not work

FEATURES:

* **New Data_source:**: menandmice_dns_zones
* data_source/dnsrecord attribute fqdn is added
* resource/dnsrecord attribute fqdn is added

* NOTES:

* resource/menandmice_ipam_record start marking `current_address` as deprecated. use `address`
* menandmice_dns_zone start marking `adintegrated` as deprecated. use `ad_integrated`
* menandmice_dns_zone start marking `displayname` as deprecated. use `display_name`
* menandmice_dns_zone start marking `dnsviewref` as deprecated. use `dns_view_ref`
* menandmice_dns_zone start marking `dnsviewrefs` as deprecated. use `dns_view_refs`
* menandmice_dns_zone start marking `dnssecsigned` as deprecated. use `dnssec_signed`

## 0.3.0

FEATURES:

* **New Resource:**: menandmice_range [GH-12]
* **New Data_source:**: menandmice_range [GH-12]

BUG FIXES:

* data_source/dhcpresrvations attribute id was set to datum
* data_source/dhcpscopes attribute id was set to datum
* data_source/dnszones attribute id was set to datum
* data_source/dnsrecord attribute id was set to datum
* data_source/ipam attribute id was set to datum

* resource/dnszone name was marker updateable but should trigger a recreated.

## 0.2.1 ( Aug 09, 2022)
