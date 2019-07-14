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

The following example is for the **AWS** provider:

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
  provider_name          = "AWS"
  aws_account_id         = "111111111111"
  vpc_id                 = "vpc-xxxxxxxxxxxxxxxxx"
  route_table_cidr_block = "10.1.0.0/16"
  container_id           = "${mongodbatlas_container.container.id}"
}

resource "aws_vpc_peering_connection_accepter" "atlas" {
  vpc_peering_connection_id = "${mongodbatlas_vpc_peering_connection.peering.connection_id}"
  auto_accept               = true
}
```

The following example is for the **GCP** provider:

```hcl
provider "mongodbatlas" {
  username = "${var.mongodb_atlas_api_username}"
  api_key = "${var.mongodb_atlas_api_key}"
}

provider "google" {
    project = "${var.gcp_project_id}"
    region = "australia-southeast1"
}

variable "gcp_project_id" {}
variable "mongodb_atlas_org_id" {}
variable "mongodb_atlas_api_username" {}
variable "mongodb_atlas_api_key" {}
variable "name_root" { default = "test" }

data "google_compute_network" "default" {
    name = "default"
}

data "mongodbatlas_container" "container" {
    group = "${mongodbatlas_project.project.id}"
    container_id = "${mongodbatlas_container.container.id}"

    depends_on = ["mongodbatlas_vpc_peering_connection.gcp_peer"]
}

resource "random_string" "name_suffix" {
    length = 6
    upper = false
    special = false
}

resource "mongodbatlas_project" "project" {
    org_id = "${var.mongodb_atlas_org_id}"
    name = "${var.name_root}-${random_string.name_suffix.result}"
}

resource "mongodbatlas_container" "container" {
    atlas_cidr_block = "192.168.100.0/18"
    provider_name = "GCP"
    group = "${mongodbatlas_project.project.id}"
    private_ip_mode = true
}

resource "mongodbatlas_vpc_peering_connection" "gcp_peer" {
    group = "${mongodbatlas_project.project.id}"
    container_id = "${mongodbatlas_container.container.id}"
    provider_name = "GCP"
    network_name = "${data.google_compute_network.default.name}"
    gcp_project_id = "${var.gcp_project_id}"
}

resource "google_compute_network_peering" "atlas_peer" {
    name = "peer-${random_string.name_suffix.result}"
    network = "${data.google_compute_network.default.self_link}"
    peer_network = "https://www.googleapis.com/compute/v1/projects/${data.mongodbatlas_container.container.gcp_project_id}/global/networks/${data.mongodbatlas_container.container.network_name}"
}
```

## Argument Reference

* `aws_account_id` ( _AWS_ ) - (Optional) AWS account ID of the owner of the peer VPC.
* `container_id` - (Required) ID of the [`mongodbatlas_container`](/docs/providers/mongodbatlas/r/container.html).

~> **NOTE:** The Atlas VPC container and the `vpc_id` peer VPC *must* share an AWS region.

* `group` - (Required) The ID of the project in which to create the VPC peering connection.
* `provider_name` - (Required) Name of the cloud provider. Valid options are:
  * `AWS`
  * `GCP`
* `route_table_cidr_block` ( _AWS_ ) - (Optional) The peer VPC CIDR block or subnet.
* `vpc_id` ( _AWS_ ) - (Optional) - The ID of the peer VPC.
* `gcp_project_id` ( _GCP_ ) - (Optional) - GCP project ID of the owner of the peer VPC.
* `network_name` ( _GCP_ ) - (Optional) - Name of the peer VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Atlas resource ID.
* `connection_id` ( _AWS_ ) - The AWS peering connection ID.
* `error_state_name` ( _AWS_ ) - The error state, if any. May be one of the following:
  * REJECTED
  * EXPIRED
  * INVALID\_ARGUMENT
* `error_message` ( _GCP_ ) - When the `status` is `FAILED` Atlas will provider a description.
* `identifier` - The same as `id`.
* `status_name` ( _AWS_ ) - Status name of the peering connection. May be one of the following:
  * INITIATING
  * PENDING\_ACCEPTANCE
  * FAILED
  * FINALIZING
  * AVAILABLE
  * TERMINATING
* `status` ( _GCP_ ) - Status of the peering connection. May be one of the following:
  * ADDING\_PEER
  * WAITING\_FOR\_USER
  * AVAILABLE
  * FAILED
  * DELETING
## Import

VPC Peering Connections can be imported using project ID and peering connection ID, in the format `PROJECTID-PEERINGID`, e.g.

```
$ terraform import mongodbatlas_vpc_peering_connection.peering 1112222b3bf99403840e8934-1aa111a1a11a111aa1a1a111
```
