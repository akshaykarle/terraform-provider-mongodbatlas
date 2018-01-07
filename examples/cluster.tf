variable "mongodb_atlas_username" {}
variable "mongodb_atlas_api_key" {}
variable "mongodb_atlas_group_id" {}

# Configure the MongoDB Atlas Provider
provider "mongodbatlas" {
  username = "${var.mongodb_atlas_username}"
  api_key = "${var.mongodb_atlas_api_key}"
}

# Create a Cluster
resource "mongodbatlas_cluster" "test" {
  name = "test"
  group = "${var.mongodb_atlas_group_id}"
  mongodb_major_version = "3.4"
  provider_name = "AWS"
  region = "US_EAST_1"
  size = "M0"
  backup = false
}
