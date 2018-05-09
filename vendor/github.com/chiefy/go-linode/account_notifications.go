package golinode

import (
	"time"

	"github.com/go-resty/resty"
)

// Notifications represent account events across all Linode things the
// account is privy to (API endpoint /account/events)
type Notification struct {
	UntilStr string `json:"until"`
	WhenStr  string `json:"when"`

	Label    string
	Message  string
	Type     string
	Severity string
	Entity   *NotificationEntity
	Until    *time.Time `json:"-"`
	When     *time.Time `json:"-"`
}

type NotificationEntity struct {
	ID    int
	Label string
	Type  string
	URL   string
}

// NotificationsPagedResponse represents a paginated Notifications API response
type NotificationsPagedResponse struct {
	*PageOptions
	Data []*Notification
}

// Endpoint gets the endpoint URL for Notification
func (NotificationsPagedResponse) Endpoint(c *Client) string {
	endpoint, err := c.Notifications.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

// Endpoint gets the endpoint URL for Notification
func (NotificationsPagedResponse) EndpointWithID(c *Client, id int) string {
	endpoint, err := c.Notifications.EndpointWithID(id)
	if err != nil {
		panic(err)
	}
	return endpoint
}

// AppendData appends Notifications when processing paginated Notification responses
func (resp *NotificationsPagedResponse) AppendData(r *NotificationsPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// SetResult sets the Resty response type of Notifications
func (NotificationsPagedResponse) SetResult(r *resty.Request) {
	r.SetResult(NotificationsPagedResponse{})
}

// ListNotifications lists Notifications
func (c *Client) ListNotifications(opts *ListOptions) ([]*Notification, error) {
	response := NotificationsPagedResponse{}
	err := c.ListHelper(response, opts)
	for _, el := range response.Data {
		el.fixDates()
	}
	return response.Data, err
}

// fixDates converts JSON timestamps to Go time.Time values
func (v *Notification) fixDates() *Notification {
	v.Until, _ = parseDates(v.UntilStr)
	v.When, _ = parseDates(v.WhenStr)
	return v
}
