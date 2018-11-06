# create vars.auto.tf with the following content:
variable "mongodb_atlas_username" {}
variable "mongodb_atlas_api_key" {}
variable "mongodb_atlas_org_id" {}
variable "pagerduty_service_key" {}

# Configure the MongoDB Atlas Provider
provider "mongodbatlas" {
  username = "${var.mongodb_atlas_username}"
  api_key = "${var.mongodb_atlas_api_key}"
}



module "cluster" {
  source = "./examples/basic"
  mongodb_atlas_org_id = "${var.mongodb_atlas_org_id}"
  project_name = "terraform"
  cluster_name = "test"
}

module "alerts" {
  source = "./examples/alerts"
  group_id = "${module.cluster.group_id}"
  pagerduty_service_key = "${var.pagerduty_service_key}"
  cluster_name = "test"
}
