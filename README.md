# Terraform Provider Menandmice

## Manual Build and Install

### Install on Mac or Linux

```shell
make install
```

### Install on Windows

First, build and install the provider.

```shell
go build -o terraform-provider-menandmice.exe
```
copy the binary: terraform-provider-menandmice.exe to:

* terraform-0.12 -> `%APPDATA%\terraform.d\plugins\windows_amd64\`
* terraform-0.14 -> `%APPDATA%\terraform.d\plugins\registry.terraform.io\local\menandmice\0.2.0\windows_amd64\`

```shell
terraform.exe init
```

# run Acceptation test

You need a working Micetro server with:
  - dnsserver: ext-master.mmdemo.net.
  - dhcpserver: DHCPScopes/192.168.2.128/25"
  - ipam-properties: location


```shell
# set provider setting that are not set in main.tf
export MENANDMICE_ENDPOINT=<api-endpoint>
export MENANDMICE_USERNAME=<your username>
export MENANDMICE_PASSWORD=<your password>

make testacc

```
