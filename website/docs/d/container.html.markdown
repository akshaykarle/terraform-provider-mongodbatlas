---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: container"
sidebar_current: "docs-mongodbatlas-datasource-container"
description: |-
    Provides details about a specific Container
---

# Data Source: mongodbatlas_container

`mongodbatlas_container` provides details about a specific Container. This represents an AWS VPC in MongoDB Atlas's network for use in VPC Peering.

This data source can prove useful when looking up the details of a previously created Container.

-> **NOTE:** Groups and projects are synonymous terms. `group` arguments on resources are the project ID.

## Example Usage

```hcl
data "mongodbatlas_project" "project" {
  name = "my-project"
}

data "mongodbatlas_container" "example" {
  group      = "${data.mongodbatlas_project.project.id}"
  identifier = "1112222b3bf99403840e8934"
}
```

## Argument Reference

* `group` - (Required) The ID of the project that the desired container belongs to.
* `identifier` - (Required) The ID of the desired container.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The same as `identifier`.
* `atlas_cidr_block` - CIDR block of the Container.
* `provider_name` - Name of provider hosting the Container. e.g. `AWS`.
* `provisioned` - Flag that indicates if the backing VPC has been created.
* `region` - Atlas-style name of region containing the Container. e.g. `US_EAST_1`
* `vpc_id` - The ID of the project's VPC. This may be empty when `provisioned` is `false`
