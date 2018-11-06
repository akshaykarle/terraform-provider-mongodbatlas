variable "aws_account_id" {}
variable "vpc_id" {}
variable "vpc_cidr_block" { default = "10.1.0.0/16" }


# Initiate a Peering connection
resource "mongodbatlas_vpc_peering_connection" "test" {
  group = "${mongodbatlas_project.test.id}"
  aws_account_id = "${var.aws_account_id}"
  vpc_id = "${var.vpc_id}"
  route_table_cidr_block = "${var.vpc_cidr_block}"
  container_id = "${mongodbatlas_container.test.id}"
}
