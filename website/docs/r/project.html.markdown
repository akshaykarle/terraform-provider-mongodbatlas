---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: project"
sidebar_current: "docs-mongodbatlas-resource-project"
description: |-
    Provides a Project resource.
---

# mongodbatlas_project

`mongodbatlas_project` provides a Project resource. This allows projects to be created.

## Example Usage

```hcl
resource "mongodbatlas_project" "project" {
  org_id = "${var.mongodb_atlas_org_id}"
  name   = "my-project"
}
```

## Argument Reference

* `name` - (Required) The name of the desired project.

~> **NOTE:** Changing `name` causes the provider to create a new project. Dependent resources may also be recreated.

* `org_id` - (Required) ID of the organization in which to create the project.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the project. Used for `group` arguments in resources.
* `cluster_count` - The number of Atlas clusters deployed in the project.
* `created` - The ISO-8601 formatted timestamp of when Atlas created the project.
