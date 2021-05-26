output "loadbalancer-ip" {
  value = local.loadbalancer_ingress.ip
}

output "loadbalancer-hostname" {
  value = local.loadbalancer_ingress.hostname
}

output "echo-service-endpoint" {
  value = "http://${local.loadbalancer_ingress.ip}/"
}
