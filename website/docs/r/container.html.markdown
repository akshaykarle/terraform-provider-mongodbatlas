---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: container"
sidebar_current: "docs-mongodbatlas-resource-container"
description: |-
    Provides a Container resource.
---

# mongodbatlas_container

`mongodbatlas_container` provides a Container resource. This represents an AWS VPC in MongoDB Atlas's network for use in VPC Peering.

~> **NOTE:** Only one Container can exist within a Project for each region. The provider allows you to define multiple container resources within the same project and region but this may lead to constant updates of the resources.

-> **NOTE:** The provider does not currently support deleting Containers due to a limitation in MongoDB Atlas's API.

-> **NOTE:** Groups and projects are synonymous terms. `group` arguments on resources are the project ID.

## Example Usage

```hcl
data "mongodbatlas_project" "project" {
  name = "my-project"
}

resource "mongodbatlas_container" "container" {
  group            = "${data.mongodbatlas_project.project.id}"
  atlas_cidr_block = "10.0.0.0/21"
  provider_name    = "AWS"
  region           = "US_EAST_1"
}

resource "mongodbatlas_cluster" "cluster" {
  group         = "${data.mongodbatlas_project.project.id}"
  provider_name = "AWS"
  region        = "US_EAST_1"
  size          = "M10"
  # ...
  depends_on = ["mongodbatlas_container.container"]
}
```

## Argument Reference

* `atlas_cidr_block` - (Required) CIDR block for the Atlas VPC in the Project region. This must be at least a /24 and at most a /21 in one of the following private networks:
  * 10.0.0.0/8
  * 172.16.0.0/12
  * 192.168.0.0/16

~> **NOTE:** The `atlas_cidr_block` value cannot be set or changed if an M10+ or VPC peering connection already exists in the Project region. To modify the CIDR block, remove all M10+ clusters and peering connections from the region, or use a new Project.

~> **NOTE:** The `atlas_cidr_block` must not overlap with containers in other regions, or with any VPC peering connections, of the Project.

-> **NOTE:** The size of the CIDR block affects the number of MongoDB nodes per container. See "Atlas CIDR Block" in the [official documentation](https://docs.atlas.mongodb.com/security-vpc-peering/)

* `group` - (Required) The ID of the project in which to create the container.
* `provider_name` - (Required) Name of the cloud provider. Currently only `AWS` is supported.
* `region` - (Required) Atlas-style name of the region in which to create the container. e.g. `US_EAST_1`. See [official documentation](https://docs.atlas.mongodb.com/reference/api/clusters-create-one/), `providerSettings.regionName`, for valid values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The container ID.
* `identifier` - The same as `id`.
* `provisioned` - Flag that indicates if the backing VPC has been created.
* `vpc_id` - The ID of the project's VPC. This will be empty when `provisioned` is `false`

## Import

There is no need to import a container. Define a resource in the same project and region to adopt or modify the existing container.
