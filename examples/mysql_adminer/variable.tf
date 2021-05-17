variable "linode_token" {
  description = "Linode APIv4 Personal Access Token"
}

variable "region" {
  default = "us-southeast"
}

variable "public_ssh_key" {
  description = "SSH Public Key Fingerprint"
  default     = "~/.ssh/id_rsa.pub"
}

variable "private_ssh_key" {
  description = "SSH Public Key Fingerprint"
  default     = "~/.ssh/id_rsa"
}

variable "mysql_db" {
  default = "cool-db"
  type = string
}

variable "mysql_user" {
  default = "testuser"
  type = string
}

variable "mysql_password" {
  default = "reallysecure!!"
  type = string
  sensitive = true
}

variable "adminer_port" {
  default = 8080
  type = number
}

