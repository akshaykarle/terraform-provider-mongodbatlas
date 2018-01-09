variable "mongodb_atlas_username" {}
variable "mongodb_atlas_api_key" {}
variable "mongodb_atlas_group_id" {}

# Configure the MongoDB Atlas Provider
provider "mongodbatlas" {
  username = "${var.mongodb_atlas_username}"
  api_key = "${var.mongodb_atlas_api_key}"
}

# Create a Container
resource "mongodbatlas_container" "test" {
  group = "${var.mongodb_atlas_group_id}"
  atlas_cidr_block = "10.0.0.0/21"
  provider_name = "AWS"
  region = "US_EAST_1"
}

# Create a Cluster
resource "mongodbatlas_cluster" "test" {
  name = "test"
  group = "${var.mongodb_atlas_group_id}"
  mongodb_major_version = "3.4"
  provider_name = "AWS"
  region = "US_EAST_1"
  size = "M10"
  backup = false
  disk_size_gb = 4.5
}
