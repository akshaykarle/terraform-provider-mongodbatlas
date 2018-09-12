module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "1.41.0"

  name                 = "test-atlas"
  cidr                 = "${var.aws_cidr}"
  public_subnets       = ["${split(",",var.public)}"]
  enable_nat_gateway   = false
  enable_dns_hostnames = true
  enable_dns_support   = true
  azs                  = ["${split(",",var.azs)}"]
}

data "aws_caller_identity" "current" {}

resource "mongodbatlas_container" "this" {
  group            = "${var.mongo_group}"
  atlas_cidr_block = "${var.mongo_cidr}"
  provider_name    = "AWS"
  region           = "${var.mongo_region}"
}

resource "mongodbatlas_vpc_peering_connection" "peering" {
  group                  = "${var.mongo_group}"
  aws_account_id         = "${data.aws_caller_identity.current.account_id}"
  vpc_id                 = "${module.vpc.vpc_id}"
  route_table_cidr_block = "${var.aws_cidr}"
  container_id           = "${mongodbatlas_container.this.id}"
}

resource "mongodbatlas_ip_whitelist" "vpc" {
  group      = "${var.mongo_group}"
  cidr_block = "${module.vpc.vpc_cidr_block}"
  comment    = "vpc cidr"
}

resource "aws_vpc_peering_connection_accepter" "atlas" {
  vpc_peering_connection_id = "${mongodbatlas_vpc_peering_connection.peering.connection_id}"
  auto_accept               = true
}

resource "aws_route" "atlas_peering" {
  route_table_id            = "${module.vpc.public_route_table_ids[0]}"
  destination_cidr_block    = "${var.mongo_cidr}"
  vpc_peering_connection_id = "${aws_vpc_peering_connection_accepter.atlas.id}"

  timeouts {
    create = "5m"
  }
}
