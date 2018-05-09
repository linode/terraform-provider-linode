package golinode

import (
	"time"

	"github.com/go-resty/resty"
)

// Events represent account events across all Linode things the
// account is privy to (API endpoint /account/events)
type Event struct {
	CreatedStr string `json:"created"`
	UpdatedStr string `json:"updated"`

	ID              int
	Status          string
	Action          string
	PercentComplete int `json:"percent_complete"`
	Rate            string
	Read            bool
	Seen            bool
	TimeRemaining   int
	Username        string
	Entity          *EventEntity
	Created         *time.Time `json:"-"`
}

type EventEntity struct {
	ID    int
	Label string
	Type  string
	URL   string
}

// EventsPagedResponse represents a paginated Events API response
type EventsPagedResponse struct {
	*PageOptions
	Data []*Event
}

// Endpoint gets the endpoint URL for Event
func (EventsPagedResponse) Endpoint(c *Client) string {
	endpoint, err := c.Events.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

// Endpoint gets the endpoint URL for Event
func (EventsPagedResponse) EndpointWithID(c *Client, id int) string {
	endpoint, err := c.Events.EndpointWithID(id)
	if err != nil {
		panic(err)
	}
	return endpoint
}

// AppendData appends Events when processing paginated Event responses
func (resp *EventsPagedResponse) AppendData(r *EventsPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// SetResult sets the Resty response type of Events
func (EventsPagedResponse) SetResult(r *resty.Request) {
	r.SetResult(EventsPagedResponse{})
}

// ListEvents lists Events
func (c *Client) ListEvents(opts *ListOptions) ([]*Event, error) {
	response := EventsPagedResponse{}
	err := c.ListHelper(response, opts)
	for _, el := range response.Data {
		el.fixDates()
	}
	return response.Data, err
}

// fixDates converts JSON timestamps to Go time.Time values
func (v *Event) fixDates() *Event {
	v.Created, _ = parseDates(v.CreatedStr)
	return v
}
