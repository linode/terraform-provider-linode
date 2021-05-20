variable "linode_token" {
  description = "Linode APIv4 Personal Access Token"
}

variable "region" {
  description = "The region to deploy the LKE cluster in."
  default = "us-east"
}


variable "echo_message" {
  description = "The message to be echoed by the echo server."
  default = "Provisioned using Terraform and LKE!"
}

variable "replica_count" {
  description = "The number of replicas of the echo server."
  default = 3
}

variable "pool_count" {
  description = "The number of instances to provision in the LKE cluster."
  default = 3
}