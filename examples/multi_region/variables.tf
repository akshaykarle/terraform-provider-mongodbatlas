variable "ssh_key" {
  description = "public_key pair data. See docs for formats: https://www.terraform.io/docs/providers/aws/r/key_pair.html"
}

variable "org_id" {
  description = "MongoDB Atlas organization ID. This must have payment info entered"
}
