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
Clone repository to: `$GOPATH/src/github.com/akshaykarle/terraform-provider-mongodbatlas`

```sh
$ mkdir -p $GOPATH/src/github.com/akshaykarle; cd $GOPATH/src/github.com/akshaykarle
$ git clone git@github.com:akshaykarle/terraform-provider-mongodbatlas
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/akshaykarle/terraform-provider-mongodbatlas
$ make build
```

## NOTE
The `mongodbatlas_project` and `mongodbatlas_container` resources do not destroy the project or container (vpc) in mongo atlas. This due to limitation of the mongo atlas API as it doesn't support deleting these resources.
