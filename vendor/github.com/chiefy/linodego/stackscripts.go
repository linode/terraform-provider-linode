package linodego

import (
	"fmt"
	"time"

	"github.com/go-resty/resty"
)

// Stackscript represents a linode stack script
type Stackscript struct {
	CreatedStr string `json:"created"`
	UpdatedStr string `json:"updated"`

	ID                int
	Username          string
	Label             string
	Description       string
	Images            []string
	DeploymentsTotal  int
	DeploymentsActive int
	IsPublic          bool
	Created           *time.Time `json:"-"`
	Updated           *time.Time `json:"-"`
	RevNote           string
	Script            string
	UserDefinedFields *map[string]string
	UserGravatarID    string
}

// StackscriptsPagedResponse represents a paginated Stackscript API response
type StackscriptsPagedResponse struct {
	*PageOptions
	Data []*Stackscript
}

// Endpoint gets the endpoint URL for Stackscript
func (StackscriptsPagedResponse) Endpoint(c *Client) string {
	endpoint, err := c.StackScripts.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

// AppendData appends Stackscripts when processing paginated Stackscript responses
func (resp *StackscriptsPagedResponse) AppendData(r *StackscriptsPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// SetResult sets the Resty response type of Stackscript
func (StackscriptsPagedResponse) SetResult(r *resty.Request) {
	r.SetResult(StackscriptsPagedResponse{})
}

// ListStackscripts lists Stackscripts
func (c *Client) ListStackscripts(opts *ListOptions) ([]*Stackscript, error) {
	response := StackscriptsPagedResponse{}
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
func (v *Stackscript) fixDates() *Stackscript {
	v.Created, _ = parseDates(v.CreatedStr)
	v.Updated, _ = parseDates(v.UpdatedStr)
	return v
}

// GetStackscript gets the Stackscript with the provided ID
func (c *Client) GetStackscript(id int) (*Stackscript, error) {
	e, err := c.StackScripts.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, id)
	r, err := c.R().SetResult(&Stackscript{}).Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*Stackscript).fixDates(), nil
}
