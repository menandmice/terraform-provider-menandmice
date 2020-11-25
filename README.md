# Terraform Provider Menandmice

## Manual Build and Install

### Install one Mac or Linux

```shell
make install
```

### Install on windows

First, build and install the provider.

```shell
go build -o terraform-provider-menandmice
```
copy the binary: terraform-provider-menandmice to:

* terraform-0.12 -> `%APPDATA%\terraform.d\plugins\windows_amd64\`
* terraform-0.13 -> `%APPDATA%\terraform.d\plugins\terraform-provider-menandmice\local\menandmice\0.2\windows_amd64\`


# run Acceptation test

You need a working man and mince server with:
  - dnsserver: mandm.example.net. mandm.example.com.
  - dhcpserver: mandm.example.net.
  - ipam-properties: location


```shell
# set provider setting that are not set in main.tf
export MENANDMICE_ENDPOINT=<api-endpoint>
export MENANDMICE_USERNAME=<your username>
export MENANDMICE_PASSWORD=<your password>

make testacc

```
