
# bet practice from https://www.terraform.io/plugin/sdkv2/best-practices/versioning

## 0.3.1 (Unreleased)

BREAKING CHANGES:

* data_source/dnsrecord name attribute can't end with a "." anymore

FEATURES:

* **New Data_source:**: menandmice_dns_zones
* data_source/dnsrecord attribute fqdn is added

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
