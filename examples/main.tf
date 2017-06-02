resource "linode_linode" "foobar" {
	image = "Ubuntu 14.04 LTS"
	kernel = "Latest 64 bit"
	name = "foobaz"
	group = "integration"
	region = "${var.region}"
	size = 1024
	private_networking = true
	ssh_key = "${var.ssh_key}"
	root_password = "${var.root_password}"

	provisioner "remote-exec" {
		inline = [
			"export PATH=$PATH:/usr/bin",
			# install nginx
			"sudo apt-get update",
			"sudo apt-get -y install nginx"
		]
	}
}