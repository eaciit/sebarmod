package sebarmod

import (
	"github.com/eaciit/toolkit"
)

type ClientStatus int

const (
	ClientInit    ClientStatus = 0
	ClientConnect              = 1
)

/*Client sebarmod client */
type Client struct {
	state ClientStatus
}

/*Status get client connection status*/
func (c *Client) Status() ClientStatus{
    return c.state
}

/*Connect connect to mod server*/
func (c *Client) Connect(server string) error {
	return nil
}

/*Call call fn on server*/
func (c *Client) Call(name string, data toolkit.M, output interface{}) error {
	return nil
}

/*IsConnected Check if client is connected */
func (c *Client) IsConnected() bool{
    connected := c.Status()==ClientConnect
    return connected
}