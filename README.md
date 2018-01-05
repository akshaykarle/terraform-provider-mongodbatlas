# terraform-provider-mongodbatlas
Terraform provider for [MongoDB Atlas](https://www.mongodb.com/cloud/atlas).

## Requirements
-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x
-	[Go](https://golang.org/doc/install) 1.9 (to build the provider plugin)

## Usage
```
provider "mongodbatlas" {
  version = "~> 0.1"
}
```

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
