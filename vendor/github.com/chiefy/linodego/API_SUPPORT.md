# API Support

## Linodes

- `/linode/instances`
  - [x] `GET`
  - [X] `POST`
- `/linode/instances/$id`
  - [x] `GET`
  - [ ] `PUT`
  - [ ] `DELETE`
- `/linode/instances/$id/boot`
  - [x] `POST`
- `/linode/instances/$id/clone`
  - [x] `POST`
- `/linode/instances/$id/kvmify`
  - [ ] `POST`
- `/linode/instances/$id/mutate`
  - [ ] `POST`
- `/linode/instances/$id/reboot`
  - [x] `POST`
- `/linode/instances/$id/rebuild`
  - [ ] `POST`
- `/linode/instances/$id/rescue`
  - [ ] `POST`
- `/linode/instances/$id/resize`
  - [x] `POST`
- `/linode/instances/$id/shutdown`
  - [x] `POST`
- `/linode/instances/$id/volumes`
  - [X] `GET`

### Backups

- `/linode/instances/$id/backups`
  - [X] `GET`
  - [ ] `POST`
- `/linode/instances/$id/backups/$id/restore`
  - [ ] `POST`
- `/linode/instances/$id/backups/cancel`
  - [ ] `POST`
- `/linode/instances/$id/backups/enable`
  - [ ] `POST`

### Configs

- `/linode/instances/$id/configs`
  - [X] `GET`
  - [ ] `POST`
- `/linode/instances/$id/configs/$id`
  - [X] `GET`
  - [ ] `PUT`
  - [ ] `DELETE`

### Disks

- `/linode/instances/$id/disks`
  - [X] `GET`
  - [ ] `POST`
- `/linode/instances/$id/disks/$id`
  - [X] `GET`
  - [ ] `PUT`
  - [ ] `POST`
  - [ ] `DELETE`
- `/linode/instances/$id/disks/$id/imagize`
  - [ ] `POST`
- `/linode/instances/$id/disks/$id/password`
  - [ ] `POST`
- `/linode/instances/$id/disks/$id/resize`
  - [ ] `POST`

### IPs

- `/linode/instances/$id/ips`
  - [ ] `GET`
  - [ ] `POST`
- `/linode/instances/$id/ips/$ip_address`
  - [ ] `GET`
  - [ ] `PUT`
  - [ ] `DELETE`
- `/linode/instances/$id/ips/sharing`
  - [ ] `POST`

### Kernels

- `/linode/kernels`
  - [X] `GET`
- `/linode/kernels/$id`
  - [X] `GET`

### StackScripts

- `/linode/stackscripts`
  - [x] `GET`
  - [ ] `POST`
- `/linode/stackscripts/$id`
  - [x] `GET`
  - [ ] `PUT`
  - [ ] `DELETE`

### Stats

- `/linode/instances/$id/stats`
  - [ ] `GET`
- `/linode/instances/$id/stats/$year/$month`
  - [ ] `GET`

### Types

- `/linode/types`
  - [X] `GET`
- `/linode/types/$id`
  - [X] `GET`

## Domains

- `/domains`
  - [ ] `GET`
  - [ ] `POST`
- `/domains/$id`
  - [ ] `GET`
  - [ ] `PUT`
  - [ ] `DELETE`
- `/domains/$id/clone`
  - [ ] `POST`
- `/domains/$id/records`
  - [ ] `GET`
  - [ ] `POST`
- `/domains/$id/records/$id`
  - [ ] `GET`
  - [ ] `PUT`
  - [ ] `DELETE`

## Longview

- `/longview/clients`
  - [ ] `GET`
  - [ ] `POST`
- `/longview/clients/$id`
  - [ ] `GET`
  - [ ] `PUT`
  - [ ] `DELETE`

### Subscriptions

- `/longview/subscriptions`
  - [ ] `GET`
- `/longview/subscriptions/$id`
  - [ ] `GET`

### NodeBalancers

- `/nodebalancers`
  - [ ] `GET`
  - [ ] `POST`
- `/nodebalancers/$id`
  - [ ] `GET`
  - [ ] `PUT`
  - [ ] `DELETE`

### Configs

- `/nodebalancers/$id/configs`
  - [ ] `GET`
  - [ ] `POST`
- `/nodebalancers/$id/configs/$id`
  - [ ] `GET`
  - [ ] `DELETE`
- `/nodebalancers/$id/configs/$id/nodes`
  - [ ] `GET`
  - [ ] `POST`
- `/nodebalancers/$id/configs/$id/nodes/$id`
  - [ ] `GET`
  - [ ] `PUT`
  - [ ] `DELETE`
- `/nodebalancers/$id/configs/$id/ssl`
  - [ ] `POST`

## Networking

- `/networking/ip-assign`
  - [ ] `POST`
- `/networking/ipv4`
  - [ ] `GET`
  - [ ] `POST`
- `/networking/ipv4/$address`
  - [ ] `GET`
  - [ ] `PUT`
  - [ ] `DELETE`

### IPv6

- `/networking/ipv6`
  - [ ] `GET`
- `/networking/ipv6/$address`
  - [ ] `GET`
  - [ ] `PUT`

## Regions

- `/regions`
  - [x] `GET`
- `/regions/$id`
  - [x] `GET`

## Support

- `/support/tickets`
  - [ ] `GET`
  - [ ] `POST`
- `/support/tickets/$id`
  - [ ] `GET`
- `/support/tickets/$id/attachments`
  - [ ] `POST`
- `/support/tickets/$id/replies`
  - [ ] `GET`
  - [ ] `POST`

## Account

### Events

- `/account/events`
  - [ ] `GET`
- `/account/events/$id`
  - [ ] `GET`
- `/account/events/$id/read`
  - [ ] `POST`
- `/account/events/$id/seen`
  - [ ] `POST`

### Invoices

- `/account/invoices/`
  - [ ] `GET`
- `/account/invoices/$id`
  - [ ] `GET`
- `/account/invoices/$id/items`
  - [ ] `GET`

### Notifications

- `/account/notifications`
  - [ ] `GET`

### OAuth Clients

- `/account/oauth-clients`
  - [ ] `GET`
  - [ ] `POST`
- `/account/oauth-clients/$id`
  - [ ] `GET`
  - [ ] `PUT`
  - [ ] `DELETE`
- `/account/oauth-clients/$id/reset_secret`
  - [ ] `POST`
- `/account/oauth-clients/$id/thumbnail`
  - [ ] `GET`
  - [ ] `PUT`

### Payments

- `/account/payments`
  - [ ] `GET`
  - [ ] `POST`
- `/account/payments/$id`
  - [ ] `GET`
- `/account/payments/paypal`
  - [ ] `GET`
- `/account/payments/paypal/execute`
  - [ ] `POST`

### Settings

- `/account/settings`
  - [ ] `GET`
  - [ ] `PUT`

### Users

- `/account/users`
  - [ ] `GET`
  - [ ] `POST`
- `/account/users/$username`
  - [ ] `GET`
  - [ ] `PUT`
  - [ ] `DELETE`
- `/account/users/$username/grants`
  - [ ] `GET`
  - [ ] `PUT`
- `/account/users/$username/password`
  - [ ] `POST`

## Images

- `/images`
  - [x] `GET`
- `/images/$id`
  - [x] `GET`
  - [ ] `PUT`
  - [ ] `DELETE`

## Volumes

- `/volumes`
  - [X] `GET`
  - [ ] `POST`
- `/volumes/$id`
  - [X] `GET`
  - [ ] `POST`
  - [ ] `DELETE`
- `/volumes/$id/attach`
  - [X] `POST`
- `/volumes/$id/clone`
  - [X] `POST`
- `/volumes/$id/detach`
  - [X] `POST`
- `/volumes/$id/resize`
  - [X] `POST`
