---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: database_user"
sidebar_current: "docs-mongodbatlas-resource-database_user"
description: |-
    Provides a Database User resource.
---

# mongodbatlas_database_user

`mongodbatlas_database_user` provides a Database User resource. This represents a database user which will be applied to all clusters within the project.

User's roles can be restricted to specific databases. If two clusters in the project have the same database name, the role will apply to both clusters and databases.

~> **NOTE:** All arguments including the password will be stored in the raw state as plain-text. [Read more about sensitive data in state](https://www.terraform.io/docs/state/sensitive-data.html)

-> **NOTE:** Groups and projects are synonymous terms. `group` arguments on resources are the project ID.

## Example Usage

```hcl
data "mongodbatlas_project" "project" {
  name = "my-project"
}

resource "mongodbatlas_database_user" "test" {
  username = "test"
  password = "initial_password"
  database = "admin"
  group    = "${data.mongodbatlas_project.project.id}"

  roles {
    name     = "read"
    database = "admin"
  }

  roles {
    name     = "readWrite"
    database = "mydatabase"
  }
}
```

## Argument Reference

* `database` - (Required) The user's authentication database. In MongoDB Atlas this is always the `admin` database.
* `group` - (Required) The ID of the project in which to create the database user.
* `password` - (Optional) User's initial password. This is required to create the user but may be removed after.

~> **NOTE:** Password may show up in logs, and it will be stored in the state file as plain-text. Password can be changed in the web interface to increase security.

* `roles` - (Required) Roles to grant on individual databases and collections. See [Roles](#roles) below for more details.
* `username` - (Required) Name of the database user.

### Roles

Block mapping a user's role to a database. A role grants actions on the given database. A role on the `admin` database can include privileges that apply to other databases.

* `name` - (Required) Name of the role to grant. See [Create a Database User](https://docs.atlas.mongodb.com/reference/api/database-users-create-a-user/) `roles.roleName` for valid values and restrictions.
* `database` - (Required) Name of database on which to grant role `name`.
* `collection` - (Optional) Collection for which the role applies. Only valid when `name` is set to `read` or `readWrite`. Role applies to all collections in the `database` if `collection` is not specified.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The database user's name.

## Import

Database users can be imported using project ID and username, in the format `PROJECTID-USERNAME`, e.g.

```
$ terraform import mongodbatlas_database_user.my_user 1112222b3bf99403840e8934-my_user
```

~> **NOTE:** Terraform will want to change the password after importing the user if a `password` argument is specified.
