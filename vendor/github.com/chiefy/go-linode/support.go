package golinode

import (
	"fmt"
	"time"

	"github.com/go-resty/resty"
)

// Ticket Statuses
const (
	TicketClosed = "closed"
	TicketOpen   = "open"
	TicketNew    = "new"
)

// Ticket represents a support ticket object
type Ticket struct {
	ID          int
	Attachments []string
	Closed      *time.Time `json:"-"`
	Description string
	Entity      *TicketEntity
	GravatarID  string
	Opened      *time.Time `json:"-"`
	OpenedBy    string
	Status      string
	Summary     string
	Updated     *time.Time `json:"-"`
	UpdatedBy   string
}

// TicketEntity refers a ticket to a specific entity
type TicketEntity struct {
	ID    int
	Label string
	Type  string
	URL   string
}

// TicketsPagedResponse represents a paginated ticket API response
type TicketsPagedResponse struct {
	*PageOptions
	Data []*Ticket
}

func (TicketsPagedResponse) Endpoint(c *Client) string {
	endpoint, err := c.Tickets.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

func (resp *TicketsPagedResponse) AppendData(r *TicketsPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

func (TicketsPagedResponse) SetResult(r *resty.Request) {
	r.SetResult(TicketsPagedResponse{})
}

// ListTypes lists support tickets
func (c *Client) ListTickets(opts *ListOptions) ([]*Ticket, error) {
	response := TicketsPagedResponse{}
	err := c.ListHelper(response, opts)
	return response.Data, err
}

// GetType gets the type with the provided ID
func (c *Client) GetTicket(id string) (*Ticket, error) {
	e, err := c.Tickets.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s", e, id)
	r, err := c.R().
		SetResult(&Ticket{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*Ticket), nil
}
