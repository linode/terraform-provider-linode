## 1.1.0 (October 31, 2018)

FEATURES:

* **New Resource** `linode_token`

* **New Data Resource** `linode_user`
* **New Data Resource** `linode_account`
* **New Data Resource** `linode_profile`

BUG FIXES:

* `linode_nodebalancer_config.check_passive` changes were not handled ([#4](https://github.com/terraform-providers/terraform-provider-template/issues/4))

## 1.0.0 (October 18, 2018)

FEATURES:

* **New Resource** `linode_instance` Initial work from @btobolaski!
* **New Resource** `linode_domain`
* **New Resource** `linode_domain_record`
* **New Resource** `linode_image` Thanks @akerl!
* **New Resource** `linode_nodebalancer`
* **New Resource** `linode_nodebalancer_config`
* **New Resource** `linode_nodebalancer_node`
* **New Resource** `linode_stackscript`
* **New Resource** `linode_sshkey`
* **New Resource** `linode_volume`

* **New Data Resource** `linode_image` Thanks @cliedeman!
* **New Data Resource** `linode_instance_type` Thanks @cliedeman!
* **New Data Resource** `linode_region` Thanks @cliedeman!
* **New Data Resource** `linode_sshkey`
