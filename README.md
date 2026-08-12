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
* terraform-1.13 -> `.terraform\providers\registry.terraform.io\menandmice\menandmice\0.4.1\windows_amd64` (relative path from the location where the `main.tf` is located along with the relevant resources)

It is possible to navigate to the examples directory and make the relevant modifications to the resources defined there to fit with the environment where this terraform module will be used. Once those changes have been made the the following statements can be executed to run terraform and apply those changes

```shell
terraform.exe init # To initialize working directory

terraform plan -out plan # To create a plan

terraform apply "plan" # To apply a plan
```

# run Acceptation test

You need a working Micetro server with:
  - dnsserver: ext-master.mmdemo.net.
  - dhcpserver: DHCPScopes/192.168.2.128/25"

The following Custom Properties have to be created beforehand as well
  * IP addresses
    * `location` of type `Text`
  * Networks
    * `location` of type `Text`
  * Zones
    * `owner` of type `Text`
    * `place` of type `Text`

```shell
# set provider setting that are not set in main.tf
export MENANDMICE_ENDPOINT=<api-endpoint> # Path to mmws such as http://localhost:8111 for local deployment
export MENANDMICE_USERNAME=<your username>
export MENANDMICE_PASSWORD=<your password>

make testacc

```
