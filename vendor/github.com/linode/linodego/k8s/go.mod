module github.com/linode/linodego/k8s

require (
	github.com/linode/linodego v0.20.1
	k8s.io/api v0.18.3
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v0.18.3
)

replace github.com/linode/linodego => ../

go 1.13
