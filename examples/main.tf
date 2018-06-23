resource "linode_nodebalancer" "kahoni-nb" {
   label = "kahoni.com"
   region = "${var.region}"
   client_conn_throttle = 0
   # group              = "kahoni"
}

resource "linode_nodebalancer_config" "kahoni-https" {
   port = 443
   nodebalancer = "${linode_nodebalancer.kahoni-nb.id}"
   protocol = "http"
   algorithm = "roundrobin"
   stickiness = "http_cookie"
   check = "http_body"
   check_interval = "90"
   check_timeout = "10"
   check_attempts = "3"
   check_path = "/test"
   check_body = "it works"
   check_passive = true
   cipher_suite = "recommended"
   # ssl_cert FUTURE pair with letsencrypt resource
   # ssl_key FUTURE pair with letsencrypt resource
   # https://opencredo.com/letsencrypt-terraform/
}

resource "linode_nodebalancer_config" "kahoni-http" {
   port = 80
   nodebalancer = "${linode_nodebalancer.kahoni-nb.id}"
   protocol = "http"
   algorithm = "roundrobin"
   stickiness = "http_cookie"
   check = "http_body"
   check_interval = "90"
   check_timeout = "10"
   check_attempts = "3"
   check_path = "/test"
   check_body = "it works"
   check_passive = true
}

resource "linode_nodebalancer_node" "kahoni-80-www" {
  # LABEL becomes kahoni-80-www
  config = "${linode_nodebalancer.default}"
  address = "${linode_instance.nginx.*.private_ipaddress}"
  weight = 50
  mode = "accept"
}

resource "linode_nodebalancer_node" "CNAME-www" {
  address = "${linode_instance.www2.ipv4_address}"
  weight = 50
}

resource "linode_domain" "kahoni-com" {
   # the default type = "master" .. call it domain_type?
   soa_email = "admin@kahoni.com"
   ttl_sec = "30"
   expire_sec = "30"
   refresh_sec = "30"
   name = "kahoni.com"
   # group              = "kahoni"
   # interesting that the bare address "@" could be set this way..
   # but terraform would have to do this behind the scenes
   # ip_address = "${linode_instance.haproxy-www.ipv4_address}"
}


resource "linode_record" "A-@" {
  domain = "${linode_domain.kahoni.id}"
  type = "A"
  name = "@"
  values = ["${linode_nodebalancer.kahoni-nb.ipv4}"]
}

resource "linode_record" "AAAA-@" {
  domain = "${linode_domain.kahoni.id}"
  type = "A"
  name = "@"
  values = ["${linode_nodebalancer.kahoni-nb.ipv6}"]
}

resource "linode_record" "CNAME-www" {
  domain = "${linode_domain.kahoni.id}"
  type = "CNAME"
  name = "www"
  values = "@" # should get auto-upgraded to array, can "value" be an alias to values?
}

resource "linode_instance" "nginx" {
  count              = 3
  image              = "linode/ubuntu18.04"
  kernel             = "linode/latest-64bit"
  name               = "kahoni-nginx-${count.index + 1}"
  # group              = "kahoni"
  region             = "${linode_nodebalancer.kahoni-nb.region}"
  instance_type               = "g6-nanode-1"
  private_networking = true
  ssh_key            = "${var.ssh_key}"
  root_password      = "${var.root_password}"

  provisioner "remote-exec" {
    inline = [
      "export PATH=$PATH:/usr/bin",

      # install nginx
      "sudo apt-get update",

      "sudo apt-get -y install nginx",
      "echo it works (node ${count.index + 1}) > /srv/www/index.html",
    ]
  }
}
