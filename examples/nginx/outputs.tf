output "pet" {
  value = random_pet.project.id
}

output "NodeBalancer-IPv4" {
  value = linode_nodebalancer.foo-nb.ipv4
}

output "NodeBalancer-IPv6" {
  value = linode_nodebalancer.foo-nb.ipv6
}

output "Nginx-Cluster" {
  value = linode_instance.nginx.*.ip_address
}

output "Nginx-Lonewolf" {
  value = linode_instance.simple.*.ip_address
}

output "Simple-IPv4" {
  value = "${data.linode_networking_ip.simple_v4.address}/${data.linode_networking_ip.simple_v4.prefix} (${data.linode_networking_ip.simple_v4.rdns})"
}

output "Simple-IPv6" {
  value = "${data.linode_networking_ip.simple_v6.address}/${data.linode_networking_ip.simple_v6.prefix} (${data.linode_networking_ip.simple_v6.rdns})"
}

