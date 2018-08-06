variable "region" {
  default = "us-central"
}

variable "ssh_key" {
  description = "SSH Public Key Fingerprint"
  default     = "~/.ssh/id_rsa.pub"
}

resource "random_pet" "project" {
  prefix    = "tf_test"
  separator = "_"
}

resource "random_string" "password" {
  length  = 32
  special = true
}

variable "nginx_count" {
  description = "The number of nginx web serving Linodes to create"
  default     = 3
}
