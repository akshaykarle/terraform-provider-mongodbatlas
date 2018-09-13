---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: vpc_peering_connection"
sidebar_current: "docs-mongodbatlas-resource-vpc_peering_connection"
description: |-
    Provides a VPC Peering Connection resource.
---

# mongodbatlas_vpc_peering_connection

`mongodbatlas_vpc_peering_connection` provides a VPC Peering Connection resource. This creates a peering request to other AWS VPCs.

Enable DNS hostnames and DNS resolution in the peer VPC. Resources can then connect to MongoDB Atlas clusters in the same region via private IP addresses using DNS names. See [Updating DNS Support](http://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/vpc-dns.html#vpc-dns-updating) for how to enable these options.

MongoDB Atlas only supports VPC peering with AWS VPCs in the same region. For multi-region clusters, you must create VPC peering connections per-region. Only MongoDB nodes in the same region can be accessed over the VPC peering connection. Nodes in remote regions are still accessed via public IPs.

Route table entries will need to be created in the AWS VPC for MongoDB clusters to be accessible via private IP. See [Updating Your Route Tables for a VPC Peering Connection](https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-routing.html) and [aws_route](/docs/providers/aws/r/route.html).

-> **NOTE:** VPC Peering is not available for M0 (Free Tier), M2, and M5 clusters.

-> **NOTE:** Groups and projects are synonymous terms. `group` arguments on resources are the project ID.

## Example Usage

```hcl
provider "aws" {
  region = "us-east-1"
}

data "mongodbatlas_project" "project" {
  name = "my-project"
}

resource "mongodbatlas_container" "container" {
  group            = "${data.mongodbatlas_project.project.id}"
  atlas_cidr_block = "10.0.0.0/21"
  provider_name    = "AWS"
  region           = "US_EAST_1"
}

resource "mongodbatlas_vpc_peering_connection" "peering" {
  group                  = "${data.mongodbatlas_project.project.id}"
  aws_account_id         = "111111111111"
  vpc_id                 = "vpc-xxxxxxxxxxxxxxxxx"
  route_table_cidr_block = "10.1.0.0/16"
  container_id           = "${mongodbatlas_container.container.id}"
}

resource "aws_vpc_peering_connection_acceptor" "atlas" {
  vpc_peering_connection_id = "${mongodbatlas_vpc_peering_connection.peering.connection_id}"
  auto_accept               = true
}
```

## Argument Reference

* `aws_account_id` - (Required) AWS account ID of the owner of the peer VPC.
* `container_id` - (Required) ID of the [`mongodbatlas_container`](/docs/providers/mongodbatlas/r/container.html).

~> **NOTE:** The Atlas VPC container and the `vpc_id` peer VPC *must* share an AWS region.

* `group` - (Required) The ID of the project in which to create the VPC peering connection.
* `route_table_cidr_block` - (Required) The peer VPC CIDR block or subnet.
* `vpc_id` - (Required) - The ID of the peer VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Atlas resource ID.
* `connection_id` - The AWS peering connection ID.
* `error_state_name` - The error state, if any. May be one of the following:
  * REJECTED
  * EXPIRED
  * INVALID\_ARGUMENT
* `identifier` - The same as `id`.
* `status_name` - Status name of the peering connection. May be one of the following:
  * INITIATING
  * PENDING\_ACCEPTANCE
  * FAILED
  * FINALIZING
  * AVAILABLE
  * TERMINATING
