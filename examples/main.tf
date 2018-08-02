variable "nginx_count" {
  default = 3
}

resource "linode_nodebalancer" "foo-nb" {
  label                = "${var.project_name}"
  region               = "${var.region}"
  client_conn_throttle = 0

  # group              = "foo"
}

resource "linode_nodebalancer_config" "foo-https" {
  port            = 443
  nodebalancer_id = "${linode_nodebalancer.foo-nb.id}"
  protocol        = "http"
  algorithm       = "roundrobin"
  stickiness      = "http_cookie"
  check           = "http_body"
  check_interval  = "90"
  check_timeout   = "10"
  check_attempts  = "3"
  check_path      = "/test"
  check_body      = "it works"
  check_passive   = true
  cipher_suite    = "recommended"

  # ssl_cert FUTURE pair with letsencrypt resource
  # ssl_key FUTURE pair with letsencrypt resource
  # https://opencredo.com/letsencrypt-terraform/
}

resource "linode_nodebalancer_config" "foo-http" {
  port            = 80
  nodebalancer_id = "${linode_nodebalancer.foo-nb.id}"
  protocol        = "http"
  algorithm       = "roundrobin"
  stickiness      = "http_cookie"
  check           = "http_body"
  check_interval  = "90"
  check_timeout   = "10"
  check_attempts  = "3"
  check_path      = "/test"
  check_body      = "it works"
  check_passive   = true
}

resource "linode_nodebalancer_node" "foo-http-www" {
  # LABEL becomes foo-80-www
  count           = "${var.nginx_count}"
  nodebalancer_id = "${linode_nodebalancer.foo-nb.id}"
  config_id       = "${linode_nodebalancer_config.foo-https.id}"
  label           = "${var.project_name}_nbnode_http_${var.nginx_count}"

  address = "${element(linode_instance.nginx.*.private_ip_address, count.index)}:80"
  weight  = 50
  mode    = "accept"
}

resource "linode_nodebalancer_node" "foo-https-www" {
  # LABEL becomes foo-80-www
  count           = "${var.nginx_count}"
  nodebalancer_id = "${linode_nodebalancer.foo-nb.id}"
  config_id       = "${linode_nodebalancer_config.foo-http.id}"
  label           = "${var.project_name}_nbnode_https_${var.nginx_count}"

  address = "${element(linode_instance.nginx.*.private_ip_address, count.index)}:80"
  weight  = 50
  mode    = "accept"
}

resource "linode_domain" "foo-com" {
  # the default type = "master" .. call it domain_type?
  soa_email   = "${var.project_name}@${substr(uuid(),0,8)}example.com"
  ttl_sec     = "30"
  expire_sec  = "30"
  refresh_sec = "30"
  domain      = "${var.project_name}${substr(uuid(),0,8)}example.com"
  domain_type = "master"

  # group              = "foo"
  # interesting that the bare address "@" could be set this way..
  # but terraform would have to do this behind the scenes
  # ip_address = "${linode_instance.haproxy-www.ipv4_address}"
}

resource "linode_domain_record" "A-root" {
  domain_id   = "${linode_domain.foo-com.id}"
  record_type = "A"
  name        = ""
  target      = "${linode_nodebalancer.foo-nb.ipv4}"
}

resource "linode_domain_record" "A-nginx" {
  count       = "${var.nginx_count}"
  domain_id   = "${linode_domain.foo-com.id}"
  name        = "${element(linode_instance.nginx.*.label, count.index)}"
  record_type = "A"
  target      = "${element(linode_instance.nginx.*.ip_address, count.index)}"
}

resource "linode_domain_record" "AAAA-root" {
  domain_id   = "${linode_domain.foo-com.id}"
  record_type = "AAAA"
  name        = ""
  target      = "${linode_nodebalancer.foo-nb.ipv6}"
}

resource "linode_domain_record" "CNAME-www" {
  domain_id   = "${linode_domain.foo-com.id}"
  record_type = "CNAME"
  name        = "www"
  target      = "${linode_domain.foo-com.domain}"
}

resource "linode_volume" "nginx-vol" {
  count     = "${var.nginx_count}"
  region    = "${linode_nodebalancer.foo-nb.region}"
  size      = 10
  linode_id = "${element(linode_instance.nginx.*.id, count.index)}"
  label     = "nginx-vol-${count.index}"

  connection {
    // user        = "root"
    // private_key = "${pathexpand(replace(var.ssh_key, ".pub", ""))}"
    host = "${element(linode_instance.nginx.*.ip_address,count.index)}"
  }

  provisioner "remote-exec" {
    inline = [
      "export PATH=$PATH:/usr/bin",
      "ls /dev/disk/by-id/",
      "sudo mkfs /dev/disk/by-id/scsi-0Linode_Volume_nginx-vol-${count.index}",
      "mkdir -p /srv/www",
      "echo /dev/disk/by-id/scsi-0Linode_Volume_nginx-vol-${count.index} /srv/www ext4 defaults 0 2 | sudo tee -a /etc/fstab",
      "mount -a",
      "echo it works node ${count.index + 1} > /srv/www/index.html",
    ]
  }
}

resource "linode_instance" "nginx" {
  count  = "${var.nginx_count}"
  image  = "linode/ubuntu18.04"
  kernel = "linode/latest-64bit"
  label  = "foo-nginx-${count.index + 1}"

  # group              = "foo"
  region             = "${linode_nodebalancer.foo-nb.region}"
  type               = "g6-nanode-1"
  private_networking = true
  ssh_key            = "${chomp(file(var.ssh_key))}"
  root_password      = "${var.root_password}"

  provisioner "remote-exec" {
    inline = [
      # install nginx
      "export PATH=$PATH:/usr/bin",

      "sudo apt-get -q update",
      "sudo apt-get -q -y install nginx",
    ]
  }
}
