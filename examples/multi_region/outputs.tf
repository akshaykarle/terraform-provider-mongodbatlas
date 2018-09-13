output "mongo_uri" {
  value = "${mongodbatlas_cluster.cluster.mongo_uri_with_options}"
}

output "bastion_ip_east" {
  value = "${module.us-east.bastion_ip}"
}

output "bastion_ip_west" {
  value = "${module.us-west.bastion_ip}"
}
