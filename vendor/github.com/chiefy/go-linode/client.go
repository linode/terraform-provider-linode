package golinode

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-resty/resty"
)

const (
	// APIHost Linode API hostname
	APIHost = "api.linode.com"
	// APIVersion Linode API version
	APIVersion = "v4"
	// APIProto connect to API with http(s)
	APIProto = "https"
	// Version of golinode
	Version = "1.0.0"
	// APIEnvVar environment var to check for API token
	APIEnvVar = "LINODE_TOKEN"
)

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

// R wraps resty's R method
func (c *Client) R() *resty.Request {
	return c.resty.R()
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

// PageOptions are the pagination parameters for List endpoints
type PageOptions struct {
	Page    int `url:"page,omitempty"`
	Pages   int `url:"pages,omitempty"`
	Results int `url:"results,omitempty"`
}

// ListOptions are the pagination and filtering (TODO) parameters for endpoints
type ListOptions struct {
	*PageOptions
	Filter string
}

// NewClient factory to create new Client struct
func NewClient(codeAPIToken *string, transport http.RoundTripper) (*Client, error) {
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
		SetHeader("User-Agent", fmt.Sprintf("go-linode %s https://github.com/chiefy/go-linode", Version))

	resources := map[string]*Resource{
		stackscriptsName:          NewResource(stackscriptsName, stackscriptsEndpoint, false),
		imagesName:                NewResource(imagesName, imagesEndpoint, false),
		instancesName:             NewResource(instancesName, instancesEndpoint, false),
		instanceDisksName:         NewResource(instanceDisksName, instanceDisksEndpoint, true),
		instanceConfigsName:       NewResource(instanceConfigsName, instanceConfigsEndpoint, true),
		instanceSnapshotsName:     NewResource(instanceSnapshotsName, instanceSnapshotsEndpoint, true),
		instanceIPsName:           NewResource(instanceIPsName, instanceIPsEndpoint, true),
		instanceVolumesName:       NewResource(instanceVolumesName, instanceVolumesEndpoint, true),
		ipaddressesName:           NewResource(ipaddressesName, ipaddressesEndpoint, false),
		ipv6poolsName:             NewResource(ipv6poolsName, ipv6poolsEndpoint, false),
		ipv6rangesName:            NewResource(ipv6rangesName, ipv6rangesEndpoint, false),
		regionsName:               NewResource(regionsName, regionsEndpoint, false),
		volumesName:               NewResource(volumesName, volumesEndpoint, false),
		kernelsName:               NewResource(kernelsName, kernelsEndpoint, false),
		typesName:                 NewResource(typesName, typesEndpoint, false),
		domainsName:               NewResource(domainsName, domainsEndpoint, false),
		domainRecordsName:         NewResource(domainRecordsName, domainRecordsEndpoint, true),
		longviewName:              NewResource(longviewName, longviewEndpoint, false),
		longviewclientsName:       NewResource(longviewclientsName, longviewclientsEndpoint, false),
		longviewsubscriptionsName: NewResource(longviewsubscriptionsName, longviewsubscriptionsEndpoint, false),
		nodebalancersName:         NewResource(nodebalancersName, nodebalancersEndpoint, false),
		nodebalancerconfigsName:   NewResource(nodebalancerconfigsName, nodebalancerconfigsEndpoint, true),
		nodebalancernodesName:     NewResource(nodebalancernodesName, nodebalancernodesEndpoint, true),
		ticketsName:               NewResource(ticketsName, ticketsEndpoint, false),
		accountName:               NewResource(accountName, accountEndpoint, false),
		invoicesName:              NewResource(invoicesName, invoicesEndpoint, false),
		invoiceItemsName:          NewResource(invoiceItemsName, invoiceItemsEndpoint, true),
		profileName:               NewResource(profileName, profileEndpoint, false),
		managedName:               NewResource(managedName, managedEndpoint, false),
	}

	return &Client{
		apiToken:  linodeAPIToken,
		resty:     restyClient,
		resources: resources,

		Images:                resources[imagesName],
		StackScripts:          resources[stackscriptsName],
		Instances:             resources[instancesName],
		Regions:               resources[regionsName],
		InstanceDisks:         resources[instanceDisksName],
		InstanceConfigs:       resources[instanceConfigsName],
		InstanceSnapshots:     resources[instanceSnapshotsName],
		InstanceIPs:           resources[instanceIPsName],
		InstanceVolumes:       resources[instanceVolumesName],
		IPAddresses:           resources[ipaddressesName],
		IPv6Pools:             resources[ipv6poolsName],
		IPv6Ranges:            resources[ipv6rangesName],
		Volumes:               resources[volumesName],
		Kernels:               resources[kernelsName],
		Types:                 resources[typesName],
		Domains:               resources[domainsName],
		Longview:              resources[longviewName],
		LongviewSubscriptions: resources[longviewsubscriptionsName],
		NodeBalancers:         resources[nodebalancersName],
		NodeBalancerConfigs:   resources[nodebalancerconfigsName],
		NodeBalancerNodes:     resources[nodebalancernodesName],
		Tickets:               resources[ticketsName],
		Account:               resources[accountName],
		Invoices:              resources[invoicesName],
		Profile:               resources[profileName],
		Managed:               resources[managedName],
	}, nil
}

type PagedResponse struct {
	ListResponse
	*PageOptions
}

type ListResponse interface {
	Endpoint(*Client) string
	AppendData(*resty.Response)
	SetResult(*resty.Request)
	ListHelper(*resty.Request, *ListOptions) error
}

// NewListOptions simplified construction of ListOptions using only
// the two writable properties, Page and Filter
func NewListOptions(Page int, Filter string) *ListOptions {
	return &ListOptions{PageOptions: &PageOptions{Page: Page}, Filter: Filter}

}

// ListHelper abstracts fetching and pagination for GET endpoints that
// do not require any Ids (top level endpoints).
// When opts (or opts.Page) is nil, all pages will be fetched and
// returned in a single (endpoint-specific)PagedResponse
// opts.results and opts.pages will be updated from the API response
func (c *Client) ListHelper(i interface{}, opts *ListOptions) error {
	req := c.R()
	if opts != nil && opts.PageOptions != nil && opts.Page > 0 {
		req.SetQueryParam("page", strconv.Itoa(opts.Page))
	}

	var (
		err     error
		pages   int
		results int
		r       *resty.Response
	)

	if opts != nil && len(opts.Filter) > 0 {
		req.SetHeader("X-Filter", opts.Filter)
	}

	switch v := i.(type) {
	case LinodeKernelsPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*LinodeKernelsPagedResponse).Pages
			results = r.Result().(*LinodeKernelsPagedResponse).Results
			v.AppendData(r.Result().(*LinodeKernelsPagedResponse))
		}
	case LinodeTypesPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*LinodeTypesPagedResponse).Pages
			results = r.Result().(*LinodeTypesPagedResponse).Results
			v.AppendData(r.Result().(*LinodeTypesPagedResponse))
		}
	case ImagesPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*ImagesPagedResponse).Pages
			results = r.Result().(*ImagesPagedResponse).Results
			v.AppendData(r.Result().(*ImagesPagedResponse))
		}
	case StackscriptsPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*StackscriptsPagedResponse).Pages
			results = r.Result().(*StackscriptsPagedResponse).Results
			v.AppendData(r.Result().(*StackscriptsPagedResponse))
		}
	case InstancesPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*InstancesPagedResponse).Pages
			results = r.Result().(*InstancesPagedResponse).Results
			v.AppendData(r.Result().(*InstancesPagedResponse))
		}
	case RegionsPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*RegionsPagedResponse).Pages
			results = r.Result().(*RegionsPagedResponse).Results
			v.AppendData(r.Result().(*RegionsPagedResponse))
		}
	case VolumesPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*VolumesPagedResponse).Pages
			results = r.Result().(*VolumesPagedResponse).Results
			v.AppendData(r.Result().(*VolumesPagedResponse))
		}
	case EventsPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*EventsPagedResponse).Pages
			results = r.Result().(*EventsPagedResponse).Results
			v.AppendData(r.Result().(*EventsPagedResponse))
		}
	case LongviewSubscriptionsPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*LongviewSubscriptionsPagedResponse).Pages
			results = r.Result().(*LongviewSubscriptionsPagedResponse).Results
			v.AppendData(r.Result().(*LongviewSubscriptionsPagedResponse))
		}
	case LongviewClientsPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*LongviewClientsPagedResponse).Pages
			results = r.Result().(*LongviewClientsPagedResponse).Results
			v.AppendData(r.Result().(*LongviewClientsPagedResponse))
		}
	case IPAddressesPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*IPAddressesPagedResponse).Pages
			results = r.Result().(*IPAddressesPagedResponse).Results
			v.AppendData(r.Result().(*IPAddressesPagedResponse))
		}
	case IPv6PoolsPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*IPv6PoolsPagedResponse).Pages
			results = r.Result().(*IPv6PoolsPagedResponse).Results
			v.AppendData(r.Result().(*IPv6PoolsPagedResponse))
		}
	case IPv6RangesPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*IPv6RangesPagedResponse).Pages
			results = r.Result().(*IPv6RangesPagedResponse).Results
			v.AppendData(r.Result().(*IPv6RangesPagedResponse))
			// @TODO consolidate this type with IPv6PoolsPagedResponse?
		}
	case TicketsPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*TicketsPagedResponse).Pages
			results = r.Result().(*TicketsPagedResponse).Results
			v.AppendData(r.Result().(*TicketsPagedResponse))
		}
	case InvoicesPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*InvoicesPagedResponse).Pages
			results = r.Result().(*InvoicesPagedResponse).Results
			v.AppendData(r.Result().(*InvoicesPagedResponse))
		}
	case NotificationsPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*NotificationsPagedResponse).Pages
			results = r.Result().(*NotificationsPagedResponse).Results
			v.AppendData(r.Result().(*NotificationsPagedResponse))
		}
	/**
	case AccountOauthClientsPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*AccountOauthClientsPagedResponse).Pages
			results = r.Result().(*AccountOauthClientsPagedResponse).Results
			v.AppendData(r.Result().(*AccountOauthClientsPagedResponse))
		}
	case AccountPaymentsPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*AccountPaymentsPagedResponse).Pages
			results = r.Result().(*AccountPaymentsPagedResponse).Results
			v.AppendData(r.Result().(*AccountPaymentsPagedResponse))
		}
	case AccountUsersPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*AccountUsersPagedResponse).Pages
			results = r.Result().(*AccountUsersPagedResponse).Results
			v.AppendData(r.Result().(*AccountUsersPagedResponse))
		}
	case ProfileAppsPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*ProfileAppsPagedResponse).Pages
			results = r.Result().(*ProfileAppsPagedResponse).Results
			v.AppendData(r.Result().(*ProfileAppsPagedResponse))
		}
	case ProfileTokensPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*ProfileTokensPagedResponse).Pages
			results = r.Result().(*ProfileTokensPagedResponse).Results
			v.AppendData(r.Result().(*ProfileTokensPagedResponse))
		}
	case ProfileWhitelistPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*ProfileWhitelistPagedResponse).Pages
			results = r.Result().(*ProfileWhitelistPagedResponse).Results
			v.AppendData(r.Result().(*ProfileWhitelistPagedResponse))
		}
	case ManagedContactsPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*ManagedContactsPagedResponse).Pages
			results = r.Result().(*ManagedContactsPagedResponse).Results
			v.AppendData(r.Result().(*ManagedContactsPagedResponse))
		}
	case ManagedCredentialsPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*ManagedCredentialsPagedResponse).Pages
			results = r.Result().(*ManagedCredentialsPagedResponse).Results
			v.AppendData(r.Result().(*ManagedCredentialsPagedResponse))
		}
	case ManagedIssuesPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*ManagedIssuesPagedResponse).Pages
			results = r.Result().(*ManagedIssuesPagedResponse).Results
			v.AppendData(r.Result().(*ManagedIssuesPagedResponse))
		}
	case ManagedLinodeSettingsPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*ManagedLinodeSettingsPagedResponse).Pages
			results = r.Result().(*ManagedLinodeSettingsPagedResponse).Results
			v.AppendData(r.Result().(*ManagedLinodeSettingsPagedResponse))
		}
	case ManagedServicesPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*ManagedServicesPagedResponse).Pages
			results = r.Result().(*ManagedServicesPagedResponse).Results
			v.AppendData(r.Result().(*ManagedServicesPagedResponse))
		}
	**/
	default:
		panic("Unknown ListHelper interface{} used")
	}

	if err != nil {
		return err
	}

	if opts == nil {
		for page := 2; page <= pages; page = page + 1 {
			c.ListHelper(i, &ListOptions{PageOptions: &PageOptions{Page: page}})
		}
	} else {
		if opts.PageOptions == nil {
			opts.PageOptions = &PageOptions{}
		}
		opts.Results = results
		opts.Pages = pages
	}

	return nil
}

// ListHelperWithID abstracts fetching and pagination for GET endpoints that
// require an Id (second level endpoints).
// When opts (or opts.Page) is nil, all pages will be fetched and
// returned in a single (endpoint-specific)PagedResponse
// opts.results and opts.pages will be updated from the API response
func (c *Client) ListHelperWithID(i interface{}, id int, opts *ListOptions) error {
	req := c.R()
	if opts != nil && opts.Page > 0 {
		req.SetQueryParam("page", strconv.Itoa(opts.Page))
	}

	var (
		err     error
		pages   int
		results int
		r       *resty.Response
	)

	if opts != nil && len(opts.Filter) > 0 {
		req.SetHeader("X-Filter", opts.Filter)
	}

	switch v := i.(type) {
	case InvoiceItemsPagedResponse:
		if r, err = req.SetResult(v).Get(v.EndpointWithID(c, id)); err == nil {
			pages = r.Result().(*InvoiceItemsPagedResponse).Pages
			results = r.Result().(*InvoiceItemsPagedResponse).Results
			v.AppendData(r.Result().(*InvoiceItemsPagedResponse))
		}
	case DomainRecordsPagedResponse:
		if r, err = req.SetResult(v).Get(v.EndpointWithID(c, id)); err == nil {
			pages = r.Result().(*DomainRecordsPagedResponse).Pages
			results = r.Result().(*DomainRecordsPagedResponse).Results
			v.AppendData(r.Result().(*DomainRecordsPagedResponse))
		}
	case InstanceSnapshotsPagedResponse:
		if r, err = req.SetResult(v).Get(v.EndpointWithID(c, id)); err == nil {
			pages = r.Result().(*InstanceSnapshotsPagedResponse).Pages
			results = r.Result().(*InstanceSnapshotsPagedResponse).Results
			v.AppendData(r.Result().(*InstanceSnapshotsPagedResponse))
		}
	case InstanceConfigsPagedResponse:
		if r, err = req.SetResult(v).Get(v.EndpointWithID(c, id)); err == nil {
			pages = r.Result().(*InstanceConfigsPagedResponse).Pages
			results = r.Result().(*InstanceConfigsPagedResponse).Results
			v.AppendData(r.Result().(*InstanceConfigsPagedResponse))
		}
	case InstanceDisksPagedResponse:
		if r, err = req.SetResult(v).Get(v.EndpointWithID(c, id)); err == nil {
			pages = r.Result().(*InstanceDisksPagedResponse).Pages
			results = r.Result().(*InstanceDisksPagedResponse).Results
			v.AppendData(r.Result().(*InstanceDisksPagedResponse))
		}
	case NodeBalancerConfigsPagedResponse:
		if r, err = req.SetResult(v).Get(v.EndpointWithID(c, id)); err == nil {
			pages = r.Result().(*NodeBalancerConfigsPagedResponse).Pages
			results = r.Result().(*NodeBalancerConfigsPagedResponse).Results
			v.AppendData(r.Result().(*NodeBalancerConfigsPagedResponse))
		}
	case InstanceVolumesPagedResponse:
		if r, err = req.SetResult(v).Get(v.EndpointWithID(c, id)); err == nil {
			pages = r.Result().(*InstanceVolumesPagedResponse).Pages
			results = r.Result().(*InstanceVolumesPagedResponse).Results
			v.AppendData(r.Result().(*InstanceVolumesPagedResponse))
		}
	/**
	case TicketAttachmentsPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*TicketAttachmentsPagedResponse).Pages
			results = r.Result().(*TicketAttachmentsPagedResponse).Results
			v.AppendData(r.Result().(*TicketAttachmentsPagedResponse))
		}
	case TicketRepliesPagedResponse:
		if r, err = req.SetResult(v).Get(v.Endpoint(c)); err == nil {
			pages = r.Result().(*TicketRepliesPagedResponse).Pages
			results = r.Result().(*TicketRepliesPagedResponse).Results
			v.AppendData(r.Result().(*TicketRepliesPagedResponse))
		}
	**/
	default:
		panic("Unknown ListHelper interface{} used")
	}

	if err != nil {
		return err
	}

	if opts == nil {
		for page := 2; page <= pages; page = page + 1 {
			c.ListHelper(i, &ListOptions{PageOptions: &PageOptions{Page: page}})
		}
	} else {
		if opts.PageOptions == nil {
			opts.PageOptions = &PageOptions{}
		}
		opts.Results = results
		opts.Pages = pages
	}

	return nil
}
