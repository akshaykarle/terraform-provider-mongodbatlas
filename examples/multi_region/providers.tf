provider "aws" {
  version = "~> 1.35"
  region  = "us-east-1"
  alias   = "east"
}

provider "aws" {
  region = "us-west-2"
  alias  = "west"
}

provider "mongodbatlas" {
  version = "~> 0.6"
}
