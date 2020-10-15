# Terraform Provider Menandmice

Run the following command to build the provider

```shell
go build -o terraform-provider-menandmice
# or
make build
```

## Test sample configuration

First, build and install the provider.

```shell
make install
```

Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
cd examples/terraform-0.12

# set provider setting that are not set in main.tf
export MENANDMICE_ENDPOINT=mandm.example.net
export MENANDMICE_USERNAME=<your username>
export MENANDMICE_PASSWORD=<your password>

terraform init && terraform apply

# destroy created resources
terraform destroy
```
