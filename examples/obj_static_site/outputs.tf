output "website-url" {
  value = "http://${linode_object_storage_bucket.website.label}.website-${linode_object_storage_bucket.website.cluster}.linodeobjects.com/"
}