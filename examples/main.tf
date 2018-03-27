variable "mongodb_atlas_username" {}
variable "mongodb_atlas_api_key" {}
variable "mongodb_atlas_org_id" {}
variable "aws_account_id" {}
variable "vpc_id" {}
variable "vpc_cidr_block" { default = "10.1.0.0/16" }
variable "whitelist_cidr_block" { default = "179.154.224.127/32" }
variable "database_user_test_password" { default = "mongodb" }

# Configure the MongoDB Atlas Provider
provider "mongodbatlas" {
  username = "${var.mongodb_atlas_username}"
  api_key = "${var.mongodb_atlas_api_key}"
}

# Create a Group
resource "mongodbatlas_project" "test" {
  org_id = "${var.mongodb_atlas_org_id}"
  name = "test"
}

# Create a Group IP Whitelist
resource "mongodbatlas_ip_whitelist" "test" {
  group = "${mongodbatlas_project.test.id}"
  cidr_block = "${var.whitelist_cidr_block}"
  comment = "test"
}

# Create a Container
resource "mongodbatlas_container" "test" {
  group = "${mongodbatlas_project.test.id}"
  atlas_cidr_block = "10.0.0.0/21"
  provider_name = "AWS"
  region = "US_EAST_1"
}

# Initiate a Peering connection
resource "mongodbatlas_vpc_peering_connection" "test" {
  group = "${mongodbatlas_project.test.id}"
  aws_account_id = "${var.aws_account_id}"
  vpc_id = "${var.vpc_id}"
  route_table_cidr_block = "${var.vpc_cidr_block}"
  container_id = "${mongodbatlas_container.test.id}"
}

# Create a Cluster
resource "mongodbatlas_cluster" "test" {
  name = "test"
  group = "${mongodbatlas_project.test.id}"
  mongodb_major_version = "3.6"
  provider_name = "TENANT"
  backing_provider = "AWS"
  region = "US_EAST_1"
  size = "M2"
  backup = false
  disk_gb_enabled = false
}

# Create a Database User
resource "mongodbatlas_database_user" "test" {
  username = "test"
  password = "${var.database_user_test_password}"
  database = "admin"
  group = "${mongodbatlas_project.test.id}"
  roles  = [
    {
      name = "read"
      database = "admin"
    }
  ]
}
