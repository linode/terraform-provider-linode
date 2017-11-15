# terraform-provider-linode

A plugin for adding a [Linode][1] provider.

[1]:https://www.linode.com

## Description

This is a custom plugin providing a Linode provider for [Terraform][2]. This is still a work in
progress. It currently only supports managing Linodes currently. In the future dns and
nodebalancers will be added. It only supports a subset of the various options that Linode offers.
In particular none of the alerting options are supported. They should be easy to add but will
likely be tedious.

[2]:https://www.terraform.io

## Requirements

* [`terraform`][2]

## Installation

1. Download the plugin from the [releases tab][3]
Both Linux and FreeBSD binaries are available.
2. Put it somewhere were it can permanently live, it doesn't need to be in your path.
3. Create or modify your `~/.terraformrc` file. You'll need at least this:

```
providers {
    linode = "terraform-provider-linode"
}
```

If you didn't add terraform-provider-linode to your path, you'll need to put the full path to the location of the plugin.

[3]:https://github.com/RainbowHackerHorse/terraform-provider-linode/releases

## Usage

### Provider Configuration

#### `linode`

```
provider "linode" {
  key = "$LINODE_API_KEY"
}
```

The provider options are:

* `key` - (Required) This is your linode api key. It will be read out of the environment variable `LINODE_API_KEY`.

### Resource Configuration

#### `linode_linode`

```
resource "linode_linode" "foobar" {
	image = "Ubuntu 14.04 LTS"
	kernel = "Latest 64 bit"
	name = "foobaz"
	group = "integration"
	region = "Dallas, TX, USA"
	size = 1024
	status = "on"
	ip_address = "8.8.8.8"
	private_networking = true
	private_ip_address = "192.168.10.50"
	ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCxtdizvJzTT38y2oXuoLUXbLUf9V0Jy9KsM0bgIvjUCSEbuLWCXKnWqgBmkv7iTKGZg3fx6JA10hiufdGHD7at5YaRUitGP2mvC2I68AYNZmLCGXh0hYMrrUB01OEXHaYhpSmXIBc9zUdTreL5CvYe3PAYzuBA0/lGFTnNsHosSd+suA4xfJWMr/Fr4/uxrpcy8N8BE16pm4kci5tcMh6rGUGtDEj6aE9k8OI4SRmSZJsNElsu/Z/K4zqCpkW/U06vOnRrE98j3NE07nxVOTqdAMZqopFiMP0MXWvd6XyS2/uKU+COLLc0+hVsgj+dVMTWfy8wZ58OJDsIKk/cI/7yF+GZz89Js+qYx7u9mNhpEgD4UrcRHpitlRgVhA8p6R4oBqb0m/rpKBd2BAFdcty3GIP9CWsARtsCbN6YDLJ1JN3xI34jSGC1ROktVHg27bEEiT5A75w3WJl96BlSo5zJsIZDTWlaqnr26YxNHba4ILdVLKigQtQpf8WFsnB9YzmDdb9K3w9szf5lAkb/SFXw+e+yPS9habkpOncL0oCsgag5wUGCEmZ7wpiY8QgARhuwsQUkxv1aUi/Nn7b7sAkKSkxtBI3LBXZ+vcUxZTH0ut4pe9rbrEed3ktAOF5FafjA1VtarPqqZ+g46xVO9llgpXcl3rVglFtXzTcUy09hGw== btobolaski@Brendans-MacBook-Pro.local"
	root_password = "terraform-test"
}
```

value                             | Type     | Forces New | Value Type | Description
--------------------------------- | -------- | ---------- | ---------- | -----------
`image`                           | Required | yes        | string     | The image to use when creating the linode. [^1]
`kernel`                          | Required | no         | string     | The kernel to start the linode with. If you can specify `"Latest 64-bit"` or `"Latest 32-bit"` for the most recent version of either that linode provices
`name`                            | Optional | no         | string     | The name of the linode
`group`                           | Optional | no         | string     | The group of the linode
`region`                          | Required | yes        | string     | The region that the linode will be created in
`size`                            | Required | yes        | int        | The amount of ram in the linode plan. i.e. 1024, 2048 or 4096
`ip_address`                      | Computed | n/a        | string     | The public ip address
`private_networking`              | Optional | sort of    | bool       | Whether or not to enable private networking. It can be enabled on an existing linode but it can't be disabled.
`private_ip_address`              | Computed | n/a        | string     | If private networking is enabled, it will be populated with the linode's private ip address
`ssh_key`                         | Required | yes        | string     | The full text of the public key to add to the root user
`root_password`                   | Required | yes        | string     | Unfortunately this is required by the linode api. You'll likely want to modify this on the server during provisioning (which won't force a new linode) and then disable password logins for ssh.
`helper_distro`                   | Optional | no         | bool       | Enable the Distro filesystem helper. Corrects fstab and inittab/upstart entries depending on the kernel you're booting. You want this unless you're providing your own kernel.
`manage_private_ip_automatically` | Optional | no         | bool       | Automatically creates network configuration files for your distro and places them into your filesystem. Will reboot your linode when enabled.

[^1]: While these technically could be modified, it requires destroying the root volume and creating a new volume which is practically the same as creating a new instance.

## Contributing

1. Fork the repo
2. Use [godep][3] to get the correct versions of the dependencies, via `godep get`
3. Make your changes
4. Apply `go fmt` to all of the files
5. Verify that the tests still pass
6. Submit a pull request

[3]:https://github.com/tools/godep

## Building
If you're unfamiliar with building go or terraform plugins, please visit the following links first:
https://golang.org/pkg/go/build/
https://golang.org/doc/install
https://www.terraform.io/guides/writing-custom-terraform-providers.html

1. Set up your GOPATH. A good default is $HOME/go. You can quickly set it by running `export GOPATH=$HOME/go`
2. cd to your new GOPATH.
3. Run `go get github.com/RainbowHackerHorse/terraform-provider-linode`
4. `cd $GOPATH/src/github.com/RainbowHackerHorse/terraform-provider-linode/bin/terraform-provider-linode`
5. `go build -o terraform-provider-linode` (To cross-compile for freebsd, run `GOOS=freebsd GOARCH=amd64 build -o terraform-provider-linode`)
6. Copy the resulting package to wherever you want on your filesystem.
7. Follow the above instructions for using it.
