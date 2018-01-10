variable "mongodb_atlas_username" {}
variable "mongodb_atlas_api_key" {}
variable "mongodb_atlas_group_id" {}
variable "aws_account_id" {}
variable "vpc_id" {}
variable "vpc_cidr_block" { default = "10.1.0.0/16" }
variable "database_user_test_password" { default = "mongodb" }

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

# Initiate a Peering connection
resource "mongodbatlas_vpc_peering_connection" "test" {
  group = "${var.mongodb_atlas_group_id}"
  aws_account_id = "${var.aws_account_id}"
  vpc_id = "${var.vpc_id}"
  route_table_cidr_block = "${var.vpc_cidr_block}"
  container_id = "${mongodbatlas_container.test.id}"
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

# Create a Database User
resource "mongodbatlas_database_user" "test" {
  username = "test"
  password = "${var.database_user_test_password}"
  database = "admin"
  group = "${var.mongodb_atlas_group_id}"
  roles  = [
    {
      name = "read"
      database = "admin"
    }
  ]
}
