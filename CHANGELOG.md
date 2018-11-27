## 1.3.0 (Unreleased)

ENHANCEMENTS:

* resource/linode_instance: Add `timeouts` support for `create`, `update`, and `delete` (defaults 10, 20, 10)
* resource/linode_image: Add `timeouts` support for `create` (defaults 20)
* resource/linode_volume: Add `timeouts` support for `create`, `update`, and `delete` (defaults 10, 20, 10)

## 1.2.0 (November 08, 2018)

ENHANCEMENTS:

* resource/linode\_instance: Add `tags` field
* resource/linode\_instance: Add `authorized_users` field and added `authorized_users` field to `disk`

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
