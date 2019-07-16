## 1.9.0 (Unreleased)
## 1.8.0 (July 16, 2019)

BACKWARDS INCOMPATIBILITIES / NOTES:

* resource/linode\_instance: `config.root_device` no longer supplies `/dev/root` when no `/dev/sda` device is present (this was an API work-around that is no longer needed) ([#10](https://github.com/terraform-providers/terraform-provider-linode/issues/10), [#18](https://github.com/terraform-providers/terraform-provider-linode/issues/18))

## 1.7.0 (July 08, 2019)

ENHANCEMENTS:

* Compatible with Terraform v0.12.0+
* Uses linodego v0.10.0
* Examples updated with new TF config syntax

BUG FIXES:

* The Linode API resizes disks by default when an instance is resized. This behavior is now accounted for -- Terraform will not resize disks unless a new size is specified in the config.
* The provider now waits for instance resizing to complete before attempting to issue disk resize jobs against an instance. This was required because new jobs issued to actively-resizing instances fail.
* Disk resizing and instance resizing are now executed in the correct order.

## 1.6.0 (April 10, 2019)

FEATURES:

* **New Resource** `linode_rdns`

* **New Data Resource** `linode_networking_ip`

ENHANCEMENTS:

* Builds now use `go mod`
* provider: Support custom `ua_prefix`
* provider: Support custom API endpoint `url`

BUG FIXES:

* Documentation and examples for `linode_domain` resource were missing required `type` field
* `linode_domain_record` field `ttl_sec` accepts `0`, as the API does (#35)
* `linode_domain` fields `ttl_sec`, `retry_sec`, `expire_sec`, and `refresh_sec` now accept `0`, as the API does

## 1.5.0 (February 06, 2019)

FEATURES:

* **New Data Resource** `linode_domain`

ENHANCEMENTS:

* Documentation has been revised with links to relevant Linode Guides & Tutorials

BUG FIXES:

* `linode_instance.alerts.0.network_out` and `linode_instance.alerts.0.transfer_quota` were not created or updating correctly ([#27](https://github.com/terraform-providers/terraform-provider-linode/issues/27))

## 1.4.0 (January 14, 2019)

BACKWARDS INCOMPATIBILITIES / NOTES:

* resource/linode\_instance: `tags` field is now a `TypeSet` instead of a `TypeList` ([#16](https://github.com/terraform-providers/terraform-provider-linode/issues/16))

ENHANCEMENTS:

* resource/linode\_domain: Add `tags` field
* resource/linode\_nodebalancer: Add `tags` field
* resource/linode\_volume: Add `tags` field

BUG FIXES:

* `linode_nodebalancer_node.label` was updated from optional to required, the Linode API has always required this field

## 1.3.0 (November 27, 2018)

ENHANCEMENTS:

* resource/linode\_instance: Add `timeouts` support for `create`, `update`, and `delete` (defaults 10, 20, 10)
* resource/linode\_image: Add `timeouts` support for `create` (defaults 20)
* resource/linode\_volume: Add `timeouts` support for `create`, `update`, and `delete` (defaults 10, 20, 10)

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
