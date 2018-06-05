variable "mongodb_atlas_username" {}
variable "mongodb_atlas_api_key" {}
# Look at https://docs.atlas.mongodb.com/reference/api/vpc-get-containers-list/ to get container IDs
variable "mongodb_atlas_container_id" {}
variable "project_name" { default = "test" }

# Configure the MongoDB Atlas Provider
provider "mongodbatlas" {
  username = "${var.mongodb_atlas_username}"
  api_key = "${var.mongodb_atlas_api_key}"
}

# Get a Group by Name
data "mongodbatlas_project" "test" {
  name = "${var.project_name}"
}

# Get a Container by Group Id and Container Id
data "mongodbatlas_container" "test" {
  group = "${data.mongodbatlas_project.test.id}"
  identifier = "${var.mongodb_atlas_container_id}"
}

output "project_id" { value = "${data.mongodbatlas_project.test.id}" }
output "project_name" { value = "${data.mongodbatlas_project.test.name}" }
output "container_cidr_block" { value = "${data.mongodbatlas_container.test.atlas_cidr_block}" }
