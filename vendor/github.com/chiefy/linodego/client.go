package linodego

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

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
	Version = "0.1.1"
	// APIEnvVar environment var to check for API token
	APIEnvVar = "LINODE_TOKEN"
	// APISecondsPerPoll how frequently to poll for new Events
	APISecondsPerPoll = 10
)

var DefaultUserAgent = fmt.Sprintf("linodego %s https://github.com/chiefy/linodego", Version)
var envDebug = false

// Client is a wrapper around the Resty client
type Client struct {
	resty     *resty.Client
	userAgent string
	resources map[string]*Resource
	debug     bool

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

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	// Wether or not we will enable Resty debugging output
	if apiDebug, ok := os.LookupEnv("LINODE_DEBUG"); ok {
		if parsed, err := strconv.ParseBool(apiDebug); err == nil {
			envDebug = parsed
			log.Println("[INFO] LINODE_DEBUG being set to", envDebug)
		} else {
			log.Println("[WARN] LINODE_DEBUG should be an integer, 0 or 1")
		}
	}

}

// SetUserAgent sets a custom user-agent for HTTP requests
func (c *Client) SetUserAgent(ua string) *Client {
	c.userAgent = ua
	c.resty.SetHeader("User-Agent", c.userAgent)

	return c
}

// R wraps resty's R method
func (c *Client) R(ctx context.Context) *resty.Request {
	return c.resty.R().
		ExpectContentType("application/json").
		SetHeader("Content-Type", "application/json").
		SetContext(ctx).
		SetError(APIError{})
}

// SetDebug sets the debug on resty's client
func (c *Client) SetDebug(debug bool) *Client {
	c.debug = debug
	c.resty.SetDebug(debug)
	return c
}

func (c *Client) SetBaseURL(url string) *Client {
	c.resty.SetHostURL(url)
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
func NewClient(hc *http.Client) (client Client) {
	restyClient := resty.NewWithClient(hc)
	client.resty = restyClient
	client.SetUserAgent(DefaultUserAgent)
	client.SetBaseURL(fmt.Sprintf("%s://%s/%s", APIProto, APIHost, APIVersion))

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
		eventsName:                NewResource(&client, eventsName, eventsEndpoint, false, Event{}, EventsPagedResponse{}),
		invoicesName:              NewResource(&client, invoicesName, invoicesEndpoint, false, Invoice{}, InvoicesPagedResponse{}),
		invoiceItemsName:          NewResource(&client, invoiceItemsName, invoiceItemsEndpoint, true, InvoiceItem{}, InvoiceItemsPagedResponse{}),
		profileName:               NewResource(&client, profileName, profileEndpoint, false, nil, nil), // really?
		managedName:               NewResource(&client, managedName, managedEndpoint, false, nil, nil), // really?
	}

	client.resources = resources

	client.SetDebug(envDebug)
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
	client.DomainRecords = resources[domainRecordsName]
	client.Longview = resources[longviewName]
	client.LongviewSubscriptions = resources[longviewsubscriptionsName]
	client.NodeBalancers = resources[nodebalancersName]
	client.NodeBalancerConfigs = resources[nodebalancerconfigsName]
	client.NodeBalancerNodes = resources[nodebalancernodesName]
	client.Tickets = resources[ticketsName]
	client.Account = resources[accountName]
	client.Events = resources[eventsName]
	client.Invoices = resources[invoicesName]
	client.Profile = resources[profileName]
	client.Managed = resources[managedName]
	return
}

// WaitForInstanceStatus waits for the Linode instance to reach the desired state
// before returning. It will timeout with an error after timeoutSeconds.
func WaitForInstanceStatus(ctx context.Context, client *Client, instanceID int, status InstanceStatus, timeoutSeconds int) error {
	start := time.Now()
	for {
		instance, err := client.GetInstance(ctx, instanceID)
		if err != nil {
			return err
		}
		complete := (instance.Status == status)

		if complete {
			return nil
		}

		time.Sleep(1 * time.Second)
		if time.Since(start) > time.Duration(timeoutSeconds)*time.Second {
			return fmt.Errorf("Instance %d didn't reach '%s' status in %d seconds", instanceID, status, timeoutSeconds)
		}
	}
}

// WaitForVolumeStatus waits for the Volume to reach the desired state
// before returning. It will timeout with an error after timeoutSeconds.
func WaitForVolumeStatus(ctx context.Context, client *Client, volumeID int, status VolumeStatus, timeoutSeconds int) error {
	start := time.Now()
	for {
		volume, err := client.GetVolume(ctx, volumeID)
		if err != nil {
			return err
		}
		complete := (volume.Status == status)

		if complete {
			return nil
		}

		time.Sleep(1 * time.Second)
		if time.Since(start) > time.Duration(timeoutSeconds)*time.Second {
			return fmt.Errorf("Volume %d didn't reach '%s' status in %d seconds", volumeID, status, timeoutSeconds)
		}
	}
}

// WaitForVolumeLinodeID waits for the Volume to match the desired LinodeID
// before returning. An active Instance will not immediately attach or detach a volume, so the
// the LinodeID must be polled to determine volume readiness from the API.
// WaitForVolumeLinodeID will timeout with an error after timeoutSeconds.
func WaitForVolumeLinodeID(ctx context.Context, client *Client, volumeID int, linodeID *int, timeoutSeconds int) error {
	start := time.Now()
	for {
		volume, err := client.GetVolume(ctx, volumeID)
		if err != nil {
			return err
		}

		if linodeID == nil && volume.LinodeID == nil {
			return nil
		} else if linodeID == nil || volume.LinodeID == nil {
			// continue waiting
		} else if *volume.LinodeID == *linodeID {
			return nil
		}

		time.Sleep(1 * time.Second)
		if time.Since(start) > time.Duration(timeoutSeconds)*time.Second {
			return fmt.Errorf("Volume %d didn't match LinodeID %d in %d seconds", volumeID, linodeID, timeoutSeconds)
		}
	}
}

// WaitForEventFinished waits for an entity action to reach the 'finished' state
// before returning. It will timeout with an error after timeoutSeconds.
// If the event indicates a failure both the failed event and the error will be returned.
func (c Client) WaitForEventFinished(ctx context.Context, id interface{}, entityType EntityType, action EventAction, minStart time.Time, timeoutSeconds int) (*Event, error) {
	start := time.Now()
	for {
		filter, err := json.Marshal(map[string]interface{}{
			// Entity is not filtered by the API
			// Perhaps one day they will permit Entity ID/Type filtering.
			// We'll have to verify these values manually, for now.
			//"entity": map[string]interface{}{
			//	"id":   fmt.Sprintf("%v", id),
			//	"type": entityType,
			//},

			// Nor is action
			//"action": action,

			// Created is not correctly filtered by the API
			// We'll have to verify these values manually, for now.
			//"created": map[string]interface{}{
			//	"+gte": minStart.Format(time.RFC3339),
			//},

			// With potentially 1000+ events coming back, we should filter on something
			"seen": false,

			// Float the latest events to page 1
			"+order_by": "created",
			"+order":    "desc",
		})

		// Optimistically restrict results to page 1.  We should remove this when more
		// precise filtering options exist.
		listOptions := NewListOptions(1, string(filter))
		events, err := c.ListEvents(ctx, listOptions)
		if err != nil {
			return nil, err
		}

		log.Printf("waiting %ds for %s events since %v for %s %v", timeoutSeconds, action, minStart, entityType, id)

		// If there are events for this instance + action, inspect them
		for _, event := range events {
			if event.Action != action {
				continue
			}
			if event.Entity.Type != entityType {
				continue
			}

			var entID string

			switch event.Entity.ID.(type) {
			case float64, float32:
				entID = fmt.Sprintf("%.f", event.Entity.ID)
			case int:
				entID = strconv.Itoa(event.Entity.ID.(int))
			default:
				entID = fmt.Sprintf("%v", event.Entity.ID)
			}

			var findID string
			switch id.(type) {
			case float64, float32:
				findID = fmt.Sprintf("%.f", id)
			case int:
				findID = strconv.Itoa(id.(int))
			default:
				findID = fmt.Sprintf("%v", id)
			}

			if entID != findID {
				// just noise..
				// log.Println(entID, "is not", id)
				continue
			} else {
				log.Println("Found event for entity.", entID, "is", id)
			}

			if *event.Created != minStart && !event.Created.After(minStart) {
				// Not the event we were looking for
				log.Println(event.Created, "is not >=", minStart)
				continue

			}

			if event.Status == EventFailed {
				return event, fmt.Errorf("%s %v action %s failed", entityType, id, action)
			} else if event.Status == EventScheduled {
				log.Printf("%s %v action %s is scheduled", entityType, id, action)
			} else if event.Status == EventFinished {
				log.Printf("%s %v action %s is finished", entityType, id, action)
				return event, nil
			}
			log.Printf("%s %v action %s is in state %s", entityType, id, action, event.Status)
		}

		// Either pushed out of the event list or hasn't been added to the list yet
		time.Sleep(time.Second * APISecondsPerPoll)
		if time.Since(start) > time.Duration(timeoutSeconds)*time.Second {
			return nil, fmt.Errorf("Did not find '%s' status of %s %v action '%s' within %d seconds", EventFinished, entityType, id, action, timeoutSeconds)
		}
	}
}
