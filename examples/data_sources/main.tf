variable "mongodb_atlas_username" {}
variable "mongodb_atlas_api_key" {}
variable "project_name" { default = "test" }

# Configure the MongoDB Atlas Provider
provider "mongodbatlas" {
  username = "${var.mongodb_atlas_username}"
  api_key = "${var.mongodb_atlas_api_key}"
}

# Create a Group
data "mongodbatlas_project" "test" {
  name = "${var.project_name}"
}

output "project_id" { value = "${data.mongodbatlas_project.test.id}" }
output "project_name" { value = "${data.mongodbatlas_project.test.name}" }
