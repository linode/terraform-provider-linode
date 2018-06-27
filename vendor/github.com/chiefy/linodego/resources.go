package linodego

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/go-resty/resty"
)

const (
	stackscriptsName          = "stackscripts"
	imagesName                = "images"
	instancesName             = "instances"
	instanceDisksName         = "disks"
	instanceConfigsName       = "configs"
	instanceIPsName           = "ips"
	instanceSnapshotsName     = "snapshots"
	instanceVolumesName       = "instancevolumes"
	ipaddressesName           = "ipaddresses"
	ipv6poolsName             = "ipv6pools"
	ipv6rangesName            = "ipv6ranges"
	regionsName               = "regions"
	volumesName               = "volumes"
	kernelsName               = "kernels"
	typesName                 = "types"
	domainsName               = "domains"
	domainRecordsName         = "records"
	longviewName              = "longview"
	longviewclientsName       = "longviewclients"
	longviewsubscriptionsName = "longviewsubscriptions"
	nodebalancersName         = "nodebalancers"
	nodebalancerconfigsName   = "nodebalancerconfigs"
	nodebalancernodesName     = "nodebalancernodes"
	ticketsName               = "tickets"
	accountName               = "account"
	eventsName                = "events"
	invoicesName              = "invoices"
	invoiceItemsName          = "invoiceitems"
	notificationsName         = "notifications"
	profileName               = "profile"
	managedName               = "managed"

	stackscriptsEndpoint          = "linode/stackscripts"
	imagesEndpoint                = "images"
	instancesEndpoint             = "linode/instances"
	instanceConfigsEndpoint       = "linode/instances/{{ .ID }}/configs"
	instanceDisksEndpoint         = "linode/instances/{{ .ID }}/disks"
	instanceSnapshotsEndpoint     = "linode/instances/{{ .ID }}/backups"
	instanceIPsEndpoint           = "linode/instances/{{ .ID }}/ips"
	instanceVolumesEndpoint       = "linode/instances/{{ .ID }}/volumes"
	ipaddressesEndpoint           = "network/ips"
	ipv6poolsEndpoint             = "network/ipv6/pools"
	ipv6rangesEndpoint            = "network/ipv6/ranges"
	regionsEndpoint               = "regions"
	volumesEndpoint               = "volumes"
	kernelsEndpoint               = "linode/kernels"
	typesEndpoint                 = "linode/types"
	domainsEndpoint               = "domains"
	domainRecordsEndpoint         = "domains/{{ .DomainID }}/records"
	longviewEndpoint              = "longview"
	longviewclientsEndpoint       = "longview/clients"
	longviewsubscriptionsEndpoint = "longview/subscriptions"
	nodebalancersEndpoint         = "nodebalancers"
	nodebalancerconfigsEndpoint   = "nodebalancers/{{ .NodeBalancerID }}/configs"
	nodebalancernodesEndpoint     = "nodebalancers/{{ .NodeBalancerID }}/configs/{{ .ConfigID }}/nodes"
	ticketsEndpoint               = "support/tickets"
	accountEndpoint               = "account"
	eventsEndpoint                = "account/events"
	invoicesEndpoint              = "account/invoices"
	invoiceItemsEndpoint          = "account/invoices/{{ .ID }}/items"
	notificationsEndpoint         = "account/notifications"
	profileEndpoint               = "profile"
	managedEndpoint               = "managed"
)

// Resource represents a linode API resource
type Resource struct {
	name             string
	endpoint         string
	isTemplate       bool
	endpointTemplate *template.Template
	R                func() *resty.Request
	PR               func() *resty.Request
}

// NewResource is the factory to create a new Resource struct. If it has a template string the useTemplate bool must be set.
func NewResource(client *Client, name string, endpoint string, useTemplate bool, singleType interface{}, pagedType interface{}) *Resource {
	var tmpl *template.Template

	if useTemplate {
		tmpl = template.Must(template.New(name).Parse(endpoint))
	}

	r := func() *resty.Request {
		return client.R().SetResult(singleType)
	}

	pr := func() *resty.Request {
		return client.R().SetResult(pagedType)
	}

	return &Resource{name, endpoint, useTemplate, tmpl, r, pr}
}

func (r Resource) render(data interface{}) (string, error) {
	if data == nil {
		return "", NewError("Cannot template endpoint with <nil> data")
	}
	out := ""
	buf := bytes.NewBufferString(out)
	if err := r.endpointTemplate.Execute(buf, struct{ ID interface{} }{data}); err != nil {
		return "", NewError(err)
	}
	return buf.String(), nil
}

// endpointWithID will return the rendered endpoint string for the resource with provided id
func (r Resource) endpointWithID(id int) (string, error) {
	if !r.isTemplate {
		return r.endpoint, nil
	}
	return r.render(id)
}

// Endpoint will return the non-templated endpoint string for resource
func (r Resource) Endpoint() (string, error) {
	if r.isTemplate {
		return "", NewError(fmt.Sprintf("Tried to get endpoint for %s without providing data for template", r.name))
	}
	return r.endpoint, nil
}
