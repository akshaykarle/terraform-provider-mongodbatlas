locals {
  shared = "${contains(list("M2", "M5"), var.cluster_tier)}"
}

# Create a Group
resource "mongodbatlas_project" "test" {
  org_id = "${var.mongodb_atlas_org_id}"
  name = "${var.project_name}"
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

# Create a Cluster
resource "mongodbatlas_cluster" "test" {
  depends_on = ["mongodbatlas_container.test"]
  name = "${var.cluster_name}"
  group = "${mongodbatlas_project.test.id}"
  mongodb_major_version = "3.6"
  provider_name = "${local.shared ? "TENANT" : "AWS"}"
  backing_provider = "${local.shared ? "AWS" : ""}"
  region = "US_EAST_1"
  size = "${var.cluster_tier}"
  backup = false
  disk_gb_enabled = "${!local.shared}"
  disk_size_gb = "${local.shared ? 0 : 10}"
}

# Create a Database User
resource "mongodbatlas_database_user" "test" {
  username = "test"
  password = "${var.database_user_test_password}"
  database = "admin"
  group = "${mongodbatlas_project.test.id}"
  roles {
    name = "read"
    database = "admin"
  }
}
