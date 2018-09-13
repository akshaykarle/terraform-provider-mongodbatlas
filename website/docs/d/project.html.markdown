---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: project"
sidebar_current: "docs-mongodbatlas-datasource-project"
description: |-
    Provides details about a specific Project
---

# Data Source: mongodbatlas_project

`mongodbatlas_project` provides details about a specific Project.

This data source can prove useful when looking up the details of a previously created Project.

~> **NOTE:** An error will be thrown if there are multiple projects with the same name in different Organizations accessible by the user.

## Example Usage

```hcl
data "mongodbatlas_project" "project" {
  name = "my-project"
}
```

## Argument Reference

* `name` - (Required) The name of the desired project.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the project. Used for `group` arguments in resources.
* `cluster_count` - The number of Atlas clusters deployed in the project.
* `created` - The ISO-8601 formatted timestamp of when Atlas created the project.
* `org_id` - The ID of the owning organization.
