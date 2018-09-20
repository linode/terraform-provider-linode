output "pet" {
  value = "${random_pet.project.id}"
}

output "NodeBalancer-IPv4" {
  value = "${linode_nodebalancer.foo-nb.0.ipv4}"
}

output "NodeBalancer-IPv6" {
  value = "${linode_nodebalancer.foo-nb.0.ipv6}"
}

output "Nginx Cluster" {
  value = "${linode_instance.nginx.*.ip_address}"
}

output "Nginx Lonewolf" {
  value = "${linode_instance.simple.*.ip_address}"
}
