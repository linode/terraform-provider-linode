resource "linode_linode" "foobar" {
  image              = "linode/ubuntu18.04"
  kernel             = "linode/latest-64bit"
  name               = "foobaz"
  # group              = "integration"
  region             = "${var.region}"
  type               = "g6-nanode-1"
  private_networking = true
  ssh_key            = "${var.ssh_key}"
  root_password      = "${var.root_password}"

  provisioner "remote-exec" {
    inline = [
      "export PATH=$PATH:/usr/bin",

      # install nginx
      "sudo apt-get update",

      "sudo apt-get -y install nginx",
    ]
  }
}
