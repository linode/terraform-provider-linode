output "Public IP" {
  value = "${linode_instance.foobar.ip_address}"
}

output "Name" {
  value = "${linode_instance.foobar.name}"
}

