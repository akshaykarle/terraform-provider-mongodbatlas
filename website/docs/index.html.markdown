---
layout: "mongodbatlas"
page_title: "Provider: mongodbatlas"
sidebar_current: "docs-mongodbatlas-index"
description: |-
  The MongoDB Atlas provider is used to interact with the resources supported by MongoDB's Atlas service. The provider needs to be configured with the proper credentials before it can be used.
---

# MongoDB Atlas Provider

The MongoDB Atlas provider is used to interact with the resources supported by
MongoDB's Atlas service. The provider needs to be configured with the proper
credentials before it can be used.

Use the nagivation to the left to read about the available resources.

-> **NOTE:** Groups and projects are synonymous terms. `group` arguments on
resources are the project ID.

## Example Usage

```hcl
# Configure the MongoDB Atlas Provider
provider "mongodbatlas" {
  username = "${var.mongodb_atlas_username}"
  api_key  = "${var.mongodb_atlas_api_key}"
}

# Create a cluster
resource "mongodbatlas_cluster" "cluster" {
  # ...
}
```

## Authentication

The MongoDB Atlas Provider can either be configured with static credentials or
environment variables for authentication. Static credentials override
environment variables.

### Static credentials

Static credentials can be provided by adding `username` and `api_key` in-line in the MongoDB Atlas provider block:

Usage:

```hcl
provider "mongodbatlas" {
  username = "username"
  api_key  = "api_key"
}
```

### Environment variables

You can provide your credentials via the `MONGODB_ATLAS_USERNAME` and
`MONGODB_ATLAS_API_KEY` environment variables:

```hcl
provider "mongodbatlas" {}
```

Usage:

```shell
$ export MONGODB_ATLAS_USERNAME="username"
$ export MONGODB_ATLAS_API_KEY="api_key"
$ terraform plan
```

## Argument Reference

In addition to [generic `provider`
arguments](https://www.terraform.io/docs/configuration/providers.html) (e.g.
`alias` and `version`), the following arguments are supported in the MongoDB
Atlas `provider` block:

* `api_key` - (Optional) This is the MongoDB Atlas API key. It must be
  provided, but it can also be sourced from the `MONGODB_ATLAS_API_KEY`
  environment variable.

* `username` - (Optional) This is the MongoDB Atlas username. It must be
  provided, but it can also be sourced from the `MONGODB_ATLAS_USERNAME`
  environment variable.
