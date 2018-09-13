---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: cluster"
sidebar_current: "docs-mongodbatlas-resource-cluster"
description: |-
    Provides a Cluster resource.
---

# mongodbatlas_cluster

`mongodbatlas_cluster` provides a Cluster resource.

-> **NOTE:** M0 (Free tier) clusters cannot be created via the API. See [Atlas M0 (Free Tier), M2, and M5 Limitations](https://docs.atlas.mongodb.com/reference/free-shared-limitations/) for more free and shared tier limitations.

-> **NOTE:** AWS users: create a [mongodbatlas_container](/docs/providers/mongodbatlas/r/container.html) in the region first if you are creating M10+ clusters and want to use VPC peering.

-> **NOTE:** The provider does not currently support [Global Clusters](https://docs.atlas.mongodb.com/global-clusters/).

-> **NOTE:** Groups and projects are synonymous terms. `group` arguments on resources are the project ID.

## Example Usage

Shared tenancy tier cluster:

```hcl
data "mongodbatlas_project" "project" {
  name = "my-project"
}

resource "mongodbatlas_cluster" "cluster" {
  name                  = "example-tenant"
  group                 = "${data.mongodbatlas_project.project.id}"
  mongodb_major_version = "3.6"
  provider_name         = "TENANT"
  backing_provider      = "AWS"
  region                = "US_EAST_1"
  size                  = "M2"
  backup                = false
  disk_gb_enabled       = false
  disk_size_gb          = 2
}
```

Single region cluster:

```hcl
resource "mongodbatlas_cluster" "cluster" {
  name                  = "example-single"
  group                 = "${data.mongodbatlas_project.project.id}"
  mongodb_major_version = "3.4"
  provider_name         = "GCP"
  region                = "CENTRAL_US"
  size                  = "M40"
  backup                = true
  disk_gb_enabled       = true
}
```

Multi-region cluster:

```hcl
resource "mongodbatlas_cluster" "cluster" {
  name                  = "example-multi"
  group                 = "${data.mongodbatlas_project.project.id}"
  mongodb_major_version = "3.4"
  provider_name         = "AZURE"
  region                = ""
  size                  = "M20"
  backup                = false
  replication_factor    = 0

  replication_spec {
    region          = "US_WEST"
    priority        = 7
    electable_nodes = 3
  }

  replication_spec {
    region          = "US_CENTRAL"
    priority        = 6
    electable_nodes = 2
  }

  replication_spec {
    region          = "US_EAST_2"
    priority        = 5
    electable_nodes = 2
    read_only_nodes = 2
  }
}
```

## Argument Reference

* `backing_provider` - (Optional) The cloud service provider for a shared tier cluster. One of `AWS`, `GCP` or `AZURE`. Only valid when `provider_name` is `TENANT`. Only `M2` and `M5` size clusters supported.
* `backup` - (Required) Enable continuous backups.
* `disk_gb_enabled` - (Optional) Enable disk auto-scaling. Defaults `true`.
* `disk_size_gb` - (Optional) AWS/GCP only. Size in GB of the server's root volume. Minimum 10. Maximum is the smaller of: instance RAM * 50 or 4096. Default value depends on instance size. See [Create a Cluster](https://docs.atlas.mongodb.com/reference/api/clusters-create-one/) `providerSettings.instanceSizeName` for default values.
* `group` - (Required) The ID of the project in which to create the cluster.
* `mongodb_major_version` - (Required) Version of the cluster to deploy. See [Create New Cluster](https://docs.atlas.mongodb.com/create-new-cluster/#select-the-mongodb-version-of-the-cluster) "Select the MongoDB Version of the Cluster" for valid versions.
* `name` - (Required) Name of the cluster.
* `num_shards` - (Optional) Set to greater than 1 to create a sharded cluster. Default 1, replica set.

-> **NOTE:** A sharded cluster cannot be converted to a replica set.

* `paused` - (Optional) Flag that indicates whether the cluster is paused. Defaults false.

-> **NOTE:** You cannot create a cluster as `paused`.

* `provider_name` - (Required) Name of the cloud provider. Current values are: `AWS`, `GCP`, `AZURE` and `TENANT`. `TENANT` also requires setting `backing_provider`.
* `region` - (Required) Atlas-style name of the region in which to create the cluster. e.g. `US_EAST_1`. See [Create a Cluster](https://docs.atlas.mongodb.com/reference/api/clusters-create-one/), `providerSettings.regionName`, for valid values. **Note:** Set to an empty string if specifying multiple `replication_spec` blocks.
* `replication_factor` - (Optional) Number of replica set members. Each shard is a replica set with the specified replication factor if a sharded cluster. Ignored if `replication_spec` is used. Possible values of 3, 5, or 7. Default 3. **Note:** Set to 0 if specifying multiple `replication_spec` blocks.
* `replication_spec` - (Optional) Configuration of each region in a multi-region cluster. See [Replication Spec](#replication-spec) below for more details.
* `size` - (Required) Instance size of all data-bearing servers in the cluster.  See [Create a Cluster](https://docs.atlas.mongodb.com/reference/api/clusters-create-one/) `providerSettings.instanceSizeName` for valid values and default resources.

### Replication Spec

The configuration of each region in a multi-region cluster.

* `electable_nodes` - (Required) Number of electable nodes to deploy to the region. Electable nodes can become the primary and can facilitate local reads. Total number of electable nodes across all regions in the cluster must be 3, 5, or 7. Specify 0 to not have electable nodes in the region. Electable nodes cannot be created if `priority` is 0.
* `priority` - (Required) Election priority of the region. Set to 0 for regions only containing read-only nodes. The first region defined with electable nodes **must** have a `priority` of 7. Following regions with electable nodes must have a priority of one less than the previous. Lowest possible priority is 1. For example, with three regions, the priorities would be: 7, 6 and 5.
* `read_only_nodes` - (Optional) Number of read-only nodes in the region. Read-only nodes can never become the primary but can facilitate local-reads. Default 0.
* `region` - (Required) Atlas-style name of the region in which to create the replica. See [Create a Cluster](https://docs.atlas.mongodb.com/reference/api/clusters-create-one/), `providerSettings.regionName`, for valid values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The container ID.
* `identifier` - The same as `id`.
* `mongodb_version` - Version of MongoDB deployed. Major.Minor.Patch.
* `mongo_uri` - Base connection string for the cluster. See `mongo_uri_with_options` for a more usable connection string.
* `mongo_uri_updated` - When the connection string was last updated. Connection string changes, for example, if you change a replica set to a sharded cluster.
* `mongo_uri_with_options` - Connection string for connecting to the Atlas cluster. Includes necessary query parameters with values appropriate for the cluster. Include a username and password for a MongoDB user associated with the project after the `mongodb://` to actually connect. See [mongodbatlas_database_user](/docs/providers/mongodbatlas/r/database_user.html) for creating users.
* `state` - Current state of the cluster. Possible states are:
  * IDLE
  * CREATING
  * UPDATING
  * DELETING
  * DELETED
  * REPAIRING

## Import

Clusters can be imported using project ID and cluster name, in the format `PROJECTID-CLUSTERNAME`, e.g.

```
$ terraform import mongodbatlas_cluster.my_cluster 1112222b3bf99403840e8934-Cluster0
```
