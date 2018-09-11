---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: ip_whitelist"
sidebar_current: "docs-mongodbatlas-resource-ip_whitelist"
description: |-
    Provides an IP Whitelist resource.
---

# mongodbatlas_ip_whitelist

`mongodbatlas_ip_whitelist` provides an IP Whitelist entry resource. The whitelist grants access from IPs or CIDRs to clusters within the Project.

-> **NOTE:** Groups and projects are synonymous terms. `group` arguments on resources are the project ID.

## Example Usage

```hcl
data "mongodbatlas_project" "project" {
  name = "my-project"
}

resource "mongodbatlas_ip_whitelist" "cidr" {
  group      = "${data.mongodbatlas_project.project.id}"
  cidr_block = "10.0.0.0/21"
  comment    = "cidr"
}

resource "mongodbatlas_ip_whitelist" "ip" {
  group      = "${data.mongodbatlas_project.project.id}"
  ip_address = "10.10.10.10"
  comment    = "ip"
}
```

## Argument Reference

* `cidr_block` - (Optional) CIDR block from which to grant access. One of `cidr_block` or `ip_address` must be specified.
* `comment` - (Optional) Comment to add to the whitelist entry.
* `group` - (Required) The ID of the project in which to add the whitelist entry.
* `ip_address` - (Optional) IP address from which to grant access. One of `cidr_block` or `ip_address` must be specified.

-> **NOTE:** The web interface allows the use of AWS security groups in the whitelist when used with VPC peering. Unfortunately there is currently a bug in the API that makes this feature incompatible with the provider. Support says they have no time frame to fix the bug as of 2018-09-12.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The container ID.

## Import

IP Whitelist entries can be imported using project ID and CIDR or IP, in the format `PROJECTID-CIDR`, e.g.

```
$ terraform import mongodbatlas_database_user.my_user 1112222b3bf99403840e8934-10.0.0.0/24
```
