package linodego

import (
	"fmt"

	"github.com/go-resty/resty"
)

// LongviewSubscription represents a LongviewSubscription object
type LongviewSubscription struct {
	ID              string
	Label           string
	ClientsIncluded int `json:"clients_included"`
	Price           *LinodePrice
	// UpdatedStr string `json:"updated"`
	// Updated *time.Time `json:"-"`
}

// LongviewSubscriptionsPagedResponse represents a paginated LongviewSubscription API response
type LongviewSubscriptionsPagedResponse struct {
	*PageOptions
	Data []*LongviewSubscription
}

// Endpoint gets the endpoint URL for LongviewSubscription
func (LongviewSubscriptionsPagedResponse) Endpoint(c *Client) string {
	endpoint, err := c.LongviewSubscriptions.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

// AppendData appends LongviewSubscriptions when processing paginated LongviewSubscription responses
func (resp *LongviewSubscriptionsPagedResponse) AppendData(r *LongviewSubscriptionsPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// SetResult sets the Resty response type of LongviewSubscription
func (LongviewSubscriptionsPagedResponse) SetResult(r *resty.Request) {
	r.SetResult(LongviewSubscriptionsPagedResponse{})
}

// ListLongviewSubscriptions lists LongviewSubscriptions
func (c *Client) ListLongviewSubscriptions(opts *ListOptions) ([]*LongviewSubscription, error) {
	response := LongviewSubscriptionsPagedResponse{}
	err := c.ListHelper(&response, opts)
	for _, el := range response.Data {
		el.fixDates()
	}
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// fixDates converts JSON timestamps to Go time.Time values
func (v *LongviewSubscription) fixDates() *LongviewSubscription {
	// v.Created, _ = parseDates(v.CreatedStr)
	// v.Updated, _ = parseDates(v.UpdatedStr)
	return v
}

// GetLongviewSubscription gets the template with the provided ID
func (c *Client) GetLongviewSubscription(id string) (*LongviewSubscription, error) {
	e, err := c.LongviewSubscriptions.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s", e, id)
	r, err := c.R().SetResult(&LongviewSubscription{}).Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*LongviewSubscription).fixDates(), nil
}
