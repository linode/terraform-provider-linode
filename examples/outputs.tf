output "Public IP" {
  value = "${linode_instance.nginx.*.ip_address}"
}

output "Name" {
  value = "${linode_instance.nginx.*.label}"
}
