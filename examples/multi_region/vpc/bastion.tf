data "aws_ami" "amzn" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-gp2"]
  }
}

resource "aws_key_pair" "ssh" {
  key_name   = "atlas_example"
  public_key = "${var.ssh_key}"
}

resource "aws_security_group" "bastion" {
  name        = "bastion_server"
  description = "Allow traffic to bastion server"
  vpc_id      = "${module.vpc.vpc_id}"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "bastion" {
  ami                    = "${data.aws_ami.amzn.id}"
  instance_type          = "t2.micro"
  key_name               = "${aws_key_pair.ssh.key_name}"
  monitoring             = false
  vpc_security_group_ids = ["${aws_security_group.bastion.id}"]
  subnet_id              = "${module.vpc.public_subnets[0]}"
  user_data              = "${file("${path.module}/cloudinit.yml")}"
}

resource "mongodbatlas_ip_whitelist" "bastion" {
  group      = "${var.mongo_group}"
  ip_address = "${aws_instance.bastion.public_ip}"
  comment    = "Bastion IP"
}
