package linodego

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-resty/resty"
)

const (
	// APIHost Linode API hostname
	APIHost = "api.linode.com"
	// APIVersion Linode API version
	APIVersion = "v4"
	// APIProto connect to API with http(s)
	APIProto = "https"
	// Version of linodego
	Version = "1.0.0"
	// APIEnvVar environment var to check for API token
	APIEnvVar = "LINODE_TOKEN"
)

var userAgent = fmt.Sprintf("linodego %s https://github.com/chiefy/linodego", Version)

// Client is a wrapper around the Resty client
type Client struct {
	apiToken  string
	resty     *resty.Client
	resources map[string]*Resource

	Images                *Resource
	InstanceDisks         *Resource
	InstanceConfigs       *Resource
	InstanceSnapshots     *Resource
	InstanceIPs           *Resource
	InstanceVolumes       *Resource
	Instances             *Resource
	IPAddresses           *Resource
	IPv6Pools             *Resource
	IPv6Ranges            *Resource
	Regions               *Resource
	StackScripts          *Resource
	Volumes               *Resource
	Kernels               *Resource
	Types                 *Resource
	Domains               *Resource
	DomainRecords         *Resource
	Longview              *Resource
	LongviewClients       *Resource
	LongviewSubscriptions *Resource
	NodeBalancers         *Resource
	NodeBalancerConfigs   *Resource
	NodeBalancerNodes     *Resource
	Tickets               *Resource
	Account               *Resource
	Invoices              *Resource
	InvoiceItems          *Resource
	Events                *Resource
	Notifications         *Resource
	Profile               *Resource
	Managed               *Resource
}

// SetUserAgent sets a custom user-agent for HTTP requests
func (c *Client) SetUserAgent(ua string) *Client {
	userAgent = ua
	return c
}

// R wraps resty's R method
func (c *Client) R() *resty.Request {
	return c.resty.R().
		SetError(APIError{})
}

// SetDebug sets the debug on resty's client
func (c *Client) SetDebug(debug bool) *Client {
	c.resty.SetDebug(debug)
	return c
}

// Resource looks up a resource by name
func (c Client) Resource(resourceName string) *Resource {
	selectedResource, ok := c.resources[resourceName]
	if !ok {
		log.Fatalf("Could not find resource named '%s', exiting.", resourceName)
	}
	return selectedResource
}

// NewClient factory to create new Client struct
func NewClient(codeAPIToken *string, transport http.RoundTripper) (client Client) {
	linodeAPIToken := ""

	if codeAPIToken != nil {
		linodeAPIToken = *codeAPIToken
	} else if envAPIToken, ok := os.LookupEnv(APIEnvVar); ok {
		linodeAPIToken = envAPIToken
	}

	if len(linodeAPIToken) == 0 || linodeAPIToken == "" {
		log.Print("Could not find LINODE_TOKEN, authenticated endpoints will fail.")
	}

	restyClient := resty.New().
		SetHostURL(fmt.Sprintf("%s://%s/%s", APIProto, APIHost, APIVersion)).
		SetAuthToken(linodeAPIToken).
		SetTransport(transport).
		SetHeader("User-Agent", userAgent)

	resources := map[string]*Resource{
		stackscriptsName:          NewResource(&client, stackscriptsName, stackscriptsEndpoint, false, Stackscript{}, StackscriptsPagedResponse{}),
		imagesName:                NewResource(&client, imagesName, imagesEndpoint, false, Image{}, ImagesPagedResponse{}),
		instancesName:             NewResource(&client, instancesName, instancesEndpoint, false, Instance{}, InstancesPagedResponse{}),
		instanceDisksName:         NewResource(&client, instanceDisksName, instanceDisksEndpoint, true, InstanceDisk{}, InstanceDisksPagedResponse{}),
		instanceConfigsName:       NewResource(&client, instanceConfigsName, instanceConfigsEndpoint, true, InstanceConfig{}, InstanceConfigsPagedResponse{}),
		instanceSnapshotsName:     NewResource(&client, instanceSnapshotsName, instanceSnapshotsEndpoint, true, InstanceSnapshot{}, InstanceSnapshotsPagedResponse{}),
		instanceIPsName:           NewResource(&client, instanceIPsName, instanceIPsEndpoint, true, InstanceIP{}, nil),                           // really?
		instanceVolumesName:       NewResource(&client, instanceVolumesName, instanceVolumesEndpoint, true, nil, InstanceVolumesPagedResponse{}), // really?
		ipaddressesName:           NewResource(&client, ipaddressesName, ipaddressesEndpoint, false, nil, IPAddressesPagedResponse{}),            // really?
		ipv6poolsName:             NewResource(&client, ipv6poolsName, ipv6poolsEndpoint, false, nil, IPv6PoolsPagedResponse{}),                  // really?
		ipv6rangesName:            NewResource(&client, ipv6rangesName, ipv6rangesEndpoint, false, IPv6Range{}, IPv6RangesPagedResponse{}),
		regionsName:               NewResource(&client, regionsName, regionsEndpoint, false, Region{}, RegionsPagedResponse{}),
		volumesName:               NewResource(&client, volumesName, volumesEndpoint, false, Volume{}, VolumesPagedResponse{}),
		kernelsName:               NewResource(&client, kernelsName, kernelsEndpoint, false, LinodeKernel{}, LinodeKernelsPagedResponse{}),
		typesName:                 NewResource(&client, typesName, typesEndpoint, false, LinodeType{}, LinodeTypesPagedResponse{}),
		domainsName:               NewResource(&client, domainsName, domainsEndpoint, false, Domain{}, DomainsPagedResponse{}),
		domainRecordsName:         NewResource(&client, domainRecordsName, domainRecordsEndpoint, true, DomainRecord{}, DomainRecordsPagedResponse{}),
		longviewName:              NewResource(&client, longviewName, longviewEndpoint, false, nil, nil), // really?
		longviewclientsName:       NewResource(&client, longviewclientsName, longviewclientsEndpoint, false, LongviewClient{}, LongviewClientsPagedResponse{}),
		longviewsubscriptionsName: NewResource(&client, longviewsubscriptionsName, longviewsubscriptionsEndpoint, false, LongviewSubscription{}, LongviewSubscriptionsPagedResponse{}),
		nodebalancersName:         NewResource(&client, nodebalancersName, nodebalancersEndpoint, false, NodeBalancer{}, NodeBalancerConfigsPagedResponse{}),
		nodebalancerconfigsName:   NewResource(&client, nodebalancerconfigsName, nodebalancerconfigsEndpoint, true, NodeBalancerConfig{}, NodeBalancerConfigsPagedResponse{}),
		nodebalancernodesName:     NewResource(&client, nodebalancernodesName, nodebalancernodesEndpoint, true, NodeBalancerNode{}, NodeBalancerNodesPagedResponse{}),
		ticketsName:               NewResource(&client, ticketsName, ticketsEndpoint, false, Ticket{}, TicketsPagedResponse{}),
		accountName:               NewResource(&client, accountName, accountEndpoint, false, Account{}, nil), // really?
		invoicesName:              NewResource(&client, invoicesName, invoicesEndpoint, false, Invoice{}, InvoicesPagedResponse{}),
		invoiceItemsName:          NewResource(&client, invoiceItemsName, invoiceItemsEndpoint, true, InvoiceItem{}, InvoiceItemsPagedResponse{}),
		profileName:               NewResource(&client, profileName, profileEndpoint, false, nil, nil), // really?
		managedName:               NewResource(&client, managedName, managedEndpoint, false, nil, nil), // really?
	}

	client.apiToken = linodeAPIToken
	client.resty = restyClient
	client.resources = resources

	client.Images = resources[imagesName]
	client.StackScripts = resources[stackscriptsName]
	client.Instances = resources[instancesName]
	client.Regions = resources[regionsName]
	client.InstanceDisks = resources[instanceDisksName]
	client.InstanceConfigs = resources[instanceConfigsName]
	client.InstanceSnapshots = resources[instanceSnapshotsName]
	client.InstanceIPs = resources[instanceIPsName]
	client.InstanceVolumes = resources[instanceVolumesName]
	client.IPAddresses = resources[ipaddressesName]
	client.IPv6Pools = resources[ipv6poolsName]
	client.IPv6Ranges = resources[ipv6rangesName]
	client.Volumes = resources[volumesName]
	client.Kernels = resources[kernelsName]
	client.Types = resources[typesName]
	client.Domains = resources[domainsName]
	client.Longview = resources[longviewName]
	client.LongviewSubscriptions = resources[longviewsubscriptionsName]
	client.NodeBalancers = resources[nodebalancersName]
	client.NodeBalancerConfigs = resources[nodebalancerconfigsName]
	client.NodeBalancerNodes = resources[nodebalancernodesName]
	client.Tickets = resources[ticketsName]
	client.Account = resources[accountName]
	client.Invoices = resources[invoicesName]
	client.Profile = resources[profileName]
	client.Managed = resources[managedName]
	return
}
