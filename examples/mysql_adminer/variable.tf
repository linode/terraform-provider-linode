variable "linode_token" {
  description = "Linode APIv4 Personal Access Token"
}

variable "region" {
  description = "The region to provision the resources in."
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
  description = "The name of the default database to be created in the MySQL server."
  default = "cool-db"
  type = string
}

variable "mysql_user" {
  description = "The name of the default user to be created in the MySQL server."
  default = "testuser"
  type = string
}

variable "mysql_password" {
  description = "The password of the default user to be created in the MySQL server."
  default = "reallysecure!!"
  type = string
  sensitive = true
}

variable "adminer_port" {
  description = "The port to expose the adminer panel on."
  default = 8080
  type = number
}

