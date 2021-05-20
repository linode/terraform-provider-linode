terraform {
  required_providers {
    linode = {
      source = "linode/linode"
    }
  }
}

provider "linode" {
  token = var.linode_token
}

data "linode_object_storage_cluster" "primary" {
  id = var.cluster
}

resource "linode_object_storage_bucket" "website" {
  cluster = data.linode_object_storage_cluster.primary.id
  label = var.bucket_name
  cors_enabled = true
  acl = "public-read"
}

resource "linode_object_storage_key" "primary" {
  label = "access-${var.bucket_name}"

  bucket_access {
    bucket_name = linode_object_storage_bucket.website.label
    cluster = linode_object_storage_bucket.website.cluster
    permissions = "read_write"
  }
}

// The output of this command is intentionally hidden in order
// to prevent the leaking of sensitive resources in stderr.
resource "null_resource" "create_website" {
  provisioner "local-exec" {
    environment = {
      ACCESS_KEY = linode_object_storage_key.primary.access_key
      SECRET_KEY = linode_object_storage_key.primary.secret_key
    }

    command = <<EOT
s3cmd ws-create \
--host-bucket="%(bucket)s.${linode_object_storage_bucket.website.cluster}.linodeobjects.com" \
--access_key=$ACCESS_KEY \
--secret_key=$SECRET_KEY \
--ws-index=${linode_object_storage_object.index.key} \
s3://${linode_object_storage_bucket.website.label} \
> /dev/null 2>&1
EOT
  }
}

resource "linode_object_storage_object" "index" {
  secret_key = linode_object_storage_key.primary.secret_key
  access_key = linode_object_storage_key.primary.access_key

  bucket = linode_object_storage_bucket.website.label
  cluster = linode_object_storage_bucket.website.cluster

  acl = "public-read"
  key = "index.html"
  source = pathexpand("./index.html")
  etag = filemd5(pathexpand("./index.html"))

  content_type = "text/html"
}

