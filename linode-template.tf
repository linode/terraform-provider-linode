provider "linode" {
  key = "your-linode-API-key-here"
}

resource "linode_linode" "your-terraform-name-here" {
        image = "Ubuntu 16.04 LTS"
        kernel = "Latest 64 bit"
        name = "your-linode-name-here"
        group = "your-linode-group-name-here"
        region = "Atlanta, GA, USA"
        size = 1024
        ssh_key = "your-ssh-id_rsa.pub-here"
        root_password = "your-root-password-here"
}
