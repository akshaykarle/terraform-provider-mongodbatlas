# terraform-provider-mongodbatlas
[![Build Status](https://travis-ci.org/akshaykarle/terraform-provider-mongodbatlas.svg?branch=master)](https://travis-ci.org/akshaykarle/terraform-provider-mongodbatlas)
[![GitHub release](https://img.shields.io/github/release/akshaykarle/terraform-provider-mongodbatlas.svg)](https://github.com/akshaykarle/terraform-provider-mongodbatlas/releases)
[![codecov](https://codecov.io/gh/akshaykarle/terraform-provider-mongodbatlas/branch/master/graph/badge.svg)](https://codecov.io/gh/akshaykarle/terraform-provider-mongodbatlas)
[![GitHub downloads](https://img.shields.io/github/downloads/akshaykarle/terraform-provider-mongodbatlas/total.svg)]()

Terraform provider for [MongoDB Atlas](https://www.mongodb.com/cloud/atlas).

## Requirements
-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x
-	[Go](https://golang.org/doc/install) 1.10 (to build the provider plugin)

## Installing the Provider
Follow the instructions to [install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin). After placing it into your plugins directory, run `terraform init` to initialize it.

## Usage
```
# Configure the MongoDB Atlas Provider
provider "mongodbatlas" {
  username = "${var.mongodb_atlas_username}"
  api_key = "${var.mongodb_atlas_api_key}"
}

# Create a Cluster
resource "mongodbatlas_cluster" "test" {
  # ...
}
```
Also look at the example under [/examples](/examples).

### Importing resources

Currently, only `mongodbatlas_cluster`, `mongodbatlas_database_user`, `mongodbatlas_vpc_peering_connection` and `mongodbatlas_ip_whitelist` can be imported.

To import any of these resources, you need the project ID (aka. group ID). This can be found in the project
settings screen.

```
# Import a cluster
terraform import mongodbatlas_cluster.example <project ID>-<cluster name>

# Import a database user
# NOTE: you'll see a plan diff for the password, this is unavoidable since the user read API omits it
terraform import mongodbatlas_database_user.example <project ID>-<username>

# Import an ip whitelist
terraform import mongodbatlas_ip_whitelist.example <project ID>-<cidr>

# Import an vpc peering
# specify the peering connection id( pcx-xxxxxxxxx )
terrform import mongodbatlas_vpc_peering_connection.example <project ID>-<pcx id>
```

## Building the Provider
Clone and build the repository

```sh
go get github.com/akshaykarle/terraform-provider-mongodbatlas
make build
```

Symlink the binary to your terraform plugins directory:

```sh
ln -s $GOPATH/bin/terraform-provider-mongodbatlas ~/.terraform.d/plugins/
```

## Updating the Provider

```sh
go get -u github.com/akshaykarle/terraform-provider-mongodbatlas
make build
```

## NOTE
The `mongodbatlas_container` resource does not destroy the container (vpc) in mongo atlas. This is due to a limitation of the mongo atlas API as it doesn't support deleting this resource.

## Contributing
* Install project dependencies: `go get github.com/kardianos/govendor`
* Run tests: `make test`
* Build the binary: `make build`
