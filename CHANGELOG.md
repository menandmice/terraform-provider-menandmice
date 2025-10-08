
# bet practice from https://www.terraform.io/plugin/sdkv2/best-practices/versioning

## 0.5.0 (Unreleased)

## 0.4.1

BUG FIXES:

* resource/range attribute `inherit_access` changed from computed to optional

BREAKING CHANGES:

* resource/range default value for attribute `subnet` changed to true


## 0.4.0

BREAKING CHANGES:

* multiple date the time format is changed to rfc 3339
  - lastmodified
  - created
  - last_seen_data
  - last_discovery_date
* data_source/dhcpscope cidr is renamed to range
* dnsrecord name attribute can't end with a "." anymore
* resouce/dnszone displayname became read only. setting it was allowed before, but did not work
* resouce/dnszone dynamic became read only. setting it was allowed before, but did not work

FEATURES:

* **New Data_source:**: menandmice_dns_zones
* **New Data_source:**: menandmice_ranges
* provider attribute `server_timezone` is added.
  Will now print the correct time for things like creation and modification dates,
  even if server is in a different time zone
* resource/dnsrecord attribute `dns_zone_ref` can now be set instead of `zone` and `server` and optional `view`
* data_source/dnsrecord attribute `fqdn` is added
* resource/dnsrecord    attribute `fqdn` is added
* data_source/range has new attribute `child_ranges`
* resource/range    has new attribute `child_ranges`
* resource/range `free_range` attribute now support searching in multiple `ranges`

NOTES:

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
