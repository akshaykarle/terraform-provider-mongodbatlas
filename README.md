# terraform-provider-mongodbatlas
[![Build Status](https://travis-ci.org/akshaykarle/terraform-provider-mongodbatlas.svg?branch=master)](https://travis-ci.org/akshaykarle/terraform-provider-mongodbatlas)
[![GitHub release](https://img.shields.io/github/release/akshaykarle/terraform-provider-mongodbatlas.svg)](https://github.com/akshaykarle/terraform-provider-mongodbatlas/releases)
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

## Building the Provider
Clone and build the repository

```sh
go get github.com/akshaykarle/terraform-provider-mongodbatlas
go build github.com/akshaykarle/terraform-provider-mongodbatlas
```

Symlink the binary to your terraform plugins directory:

```sh
ln -s $GOPATH/bin/terraform-provider-mongodbatlas ~/.terraform.d/plugins/
```

## Updating the Provider

```sh
go get -u github.com/akshaykarle/terraform-provider-mongodbatlas
go build github.com/akshaykarle/terraform-provider-mongodbatlas
```

## NOTE
The `mongodbatlas_project` and `mongodbatlas_container` resources do not destroy the project or container (vpc) in mongo atlas. This due to limitation of the mongo atlas API as it doesn't support deleting these resources.
