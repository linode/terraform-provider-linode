package linodego

import (
	"fmt"

	"github.com/go-resty/resty"
)

// DomainRecord represents a DomainRecord object
type DomainRecord struct {
	ID       int
	Type     string
	Name     string
	Target   string
	Priority int
	Weight   int
	Port     int
	Service  string
	Protocol string
	TTLSec   int `json:"ttl_sec"`
	Tag      string
}

// DomainRecordsPagedResponse represents a paginated DomainRecord API response
type DomainRecordsPagedResponse struct {
	*PageOptions
	Data []*DomainRecord
}

// endpoint gets the endpoint URL for InstanceConfig
func (DomainRecordsPagedResponse) endpointWithID(c *Client, id int) string {
	endpoint, err := c.DomainRecords.endpointWithID(id)
	if err != nil {
		panic(err)
	}
	return endpoint
}

// appendData appends DomainRecords when processing paginated DomainRecord responses
func (resp *DomainRecordsPagedResponse) appendData(r *DomainRecordsPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// setResult sets the Resty response type of DomainRecord
func (DomainRecordsPagedResponse) setResult(r *resty.Request) {
	r.SetResult(DomainRecordsPagedResponse{})
}

// ListDomainRecords lists DomainRecords
func (c *Client) ListDomainRecords(opts *ListOptions) ([]*DomainRecord, error) {
	response := DomainRecordsPagedResponse{}
	err := c.listHelper(&response, opts)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// GetDomainRecord gets the template with the provided ID
func (c *Client) GetDomainRecord(id string) (*DomainRecord, error) {
	e, err := c.DomainRecords.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s", e, id)
	r, err := c.R().SetResult(&DomainRecord{}).Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*DomainRecord), nil
}
