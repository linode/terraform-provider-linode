terraform {
  required_providers {
    linode = {
      source = "linode/linode"
    }
  }
}

provider "linode" {
  token = var.linode_token
}

data "linode_instance_type" "default" {
  id = "g6-standard-2"
}

resource "linode_firewall" "firewall" {
  label = "my-db-firewall"
  inbound_policy = "DROP"
  outbound_policy = "DROP"

  inbound {
    label = "allow-http"
    action = "ACCEPT"
    protocol = "TCP"
    ports = "${var.adminer_port}"
    ipv4 = ["0.0.0.0/0"]
  }

  linodes = [
    linode_instance.mysql.id,
    linode_instance.adminer.id
  ]
}

resource "linode_instance" "mysql" {
  label = "my-mysql-server"
  type = data.linode_instance_type.default.id
  region = var.region
  image = "linode/alpine3.13"
  authorized_keys = [
    chomp(file(var.public_ssh_key))
  ]

  interface {
    purpose = "public"
  }

  interface {
    purpose = "vlan"
    label = "my-db-vlan"
    ipam_address = "10.0.0.1/24"
  }

  connection {
    host        = linode_instance.mysql.ip_address
    type        = "ssh"
    user        = "root"
    agent       = "false"
    private_key = chomp(file(var.private_ssh_key))
  }

  provisioner "remote-exec" {
    inline = [
      "apk add docker docker-compose",
      "rc-update add docker boot",
      "service docker start",
      "timeout 15 sh -c \"until docker info; do echo .; sleep 1; done\"",
      "docker volume create mysql_data",

      "docker run -idt --restart=unless-stopped \\",
      "-p 10.0.0.1:3306:3306 \\",
      "-v mysql_data:/var/lib/mysql \\",
      "-e MYSQL_RANDOM_ROOT_PASSWORD=1 \\",
      "-e MYSQL_DATABASE=${var.mysql_db} \\",
      "-e MYSQL_USER=${var.mysql_user} \\",
      "-e MYSQL_PASSWORD=${var.mysql_password} \\",
      "mysql"
    ]
  }
}

resource "linode_instance" "adminer" {
  label = "my-adminer-server"
  type = data.linode_instance_type.default.id
  region = var.region
  image = "linode/alpine3.13"
  authorized_keys = [
    chomp(file(var.public_ssh_key))
  ]

  interface {
    purpose = "public"
  }

  interface {
    purpose = "vlan"
    label = "my-db-vlan"
    ipam_address = "10.0.0.2/24"
  }

  connection {
    host        = linode_instance.adminer.ip_address
    type        = "ssh"
    user        = "root"
    agent       = "false"
    private_key = chomp(file(var.private_ssh_key))
  }

  provisioner "remote-exec" {
    inline = [
      "apk add docker docker-compose",
      "rc-update add docker boot",
      "service docker start",
      "timeout 15 sh -c \"until docker info; do echo .; sleep 1; done\"",

      "docker run -idt --restart=unless-stopped \\",
      "-p ${var.adminer_port}:8080 \\",
      "-e ADMINER_DEFAULT_SERVER=10.0.0.1 \\",
      "adminer"
    ]
  }
}