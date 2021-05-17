output "adminer-url" {
  value = "http://${linode_instance.adminer.ip_address}:${var.adminer_port}/"
}

output "mysql-user" {
  value = var.mysql_user
}

output "mysql-vlan-ipv4" {
  value = replace(linode_instance.mysql.interface.1.ipam_address, "/24", "")
}

output "adminer-vlan-ipv4" {
  value = replace(linode_instance.adminer.interface.1.ipam_address, "/24", "")
}