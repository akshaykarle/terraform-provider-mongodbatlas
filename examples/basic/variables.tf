variable "cluster_name" {
  default = "test"
}
variable "cluster_tier" {
  default = "M2"
}
variable "database_user_test_password" { default = "mongodb" }
variable "mongodb_atlas_org_id" {}
variable "project_name" {
  description = "Name of project in MongoDB Atlas"
  default = "test"
}
variable "whitelist_cidr_block" { default = "179.154.224.127/32" }
