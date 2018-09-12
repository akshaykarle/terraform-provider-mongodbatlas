variable "aws_cidr" {
  description = "CIDR for the AWS VPC"
}

variable "public" {
  description = "Comma separated list of CIDRs for AWS Public subnets. Must be a subset of aws_cidr"
}

variable "azs" {
  description = "Comma separated list of Availability Zones for AWS VPC"
}

variable "ssh_key" {
  description = "public_key pair data. See docs for formats: https://www.terraform.io/docs/providers/aws/r/key_pair.html"
}

variable "mongo_cidr" {
  description = "CIDR for the MongoDB Atlas VPC"
}

variable "mongo_group" {
  description = "Group or Project ID for MongoDB Atlas resources"
}

variable "mongo_region" {
  description = "Region for MongoDB resources. This must match the AWS region"
}
