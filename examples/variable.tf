variable "region" {
  default = "us-central"
}

variable "root_password" {}

variable "ssh_key" {
  description = "SSH Public Key Fingerprint"
  default     = "~/.ssh/id_rsa.pub"
}

variable "project_name" {
  description = "A name for this example project.  This will be used in domain names and labels."
  default     = "tf_test_foobar"
}
