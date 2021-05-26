variable "linode_token" {
  description = "Linode APIv4 Personal Access Token"
}

variable "cluster" {
  description = "The name of the Object Storage cluster to deploy to."
  default = "us-east-1"
}

variable "bucket_name" {
  description = "The name of the Object Storage bucket to create."
  type = string
}