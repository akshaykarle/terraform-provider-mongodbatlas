Larger worked example configuration creating a dual region cluster with VPC peering connections to AWS VPCs. Each VPC has a small bastion server running Amazon Linux 2 for testing connections to the MongoDB cluster.

Terraform will output the bastion IPs and cluster connection string. Connect to a bastion server as the user `ec2-user` using SSH key given as input variable. Run `mongo -u test -p super_password <cluster_connection_string>` to connect to the cluster.
