resource "mongodbatlas_project" "project" {
  org_id = "${var.org_id}"
  name   = "example"
}

module "us-east" {
  source = "./vpc"

  providers = {
    aws = "aws.east"
  }

  aws_cidr = "10.0.0.0/21"
  public   = "10.0.0.0/24,10.0.1.0/24"
  azs      = "us-east-1b,us-east-1c"
  ssh_key  = "${var.ssh_key}"

  mongo_cidr   = "10.0.8.0/21"
  mongo_group  = "${mongodbatlas_project.project.id}"
  mongo_region = "US_EAST_1"
}

module "us-west" {
  source = "./vpc"

  providers = {
    aws = "aws.west"
  }

  aws_cidr = "10.1.0.0/21"
  public   = "10.1.0.0/24,10.1.1.0/24"
  azs      = "us-west-2b,us-west-2c"
  ssh_key  = "${var.ssh_key}"

  mongo_cidr   = "10.1.8.0/21"
  mongo_group  = "${mongodbatlas_project.project.id}"
  mongo_region = "US_WEST_2"
}

resource "mongodbatlas_cluster" "cluster" {
  name                  = "example"
  group                 = "${mongodbatlas_project.project.id}"
  mongodb_major_version = "4.0"
  provider_name         = "AWS"
  region                = ""
  size                  = "M10"
  backup                = false
  disk_gb_enabled       = false
  replication_factor    = 0

  replication_spec {
    region          = "US_EAST_1"
    priority        = 7
    electable_nodes = 3
  }

  replication_spec {
    region          = "US_WEST_2"
    priority        = 6
    electable_nodes = 2
  }

  # Atlas Containers need configuring before creating the Cluster
  depends_on = ["module.us-east", "module.us-west"]
}

resource "mongodbatlas_database_user" "test" {
  username = "test"
  database = "admin"
  password = "super_password"
  group    = "${mongodbatlas_project.project.id}"

  roles {
    name     = "atlasAdmin"
    database = "admin"
  }
}
