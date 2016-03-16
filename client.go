package sebarmod

import (
	"github.com/eaciit/toolkit"
    "errors"
)

type ClientStatus int

const (
	ClientInit    ClientStatus = 0
	ClientConnect              = 1
)

func NewClient(host string, config toolkit.M) *Client{
    c := new(Client)
    c.Host = host
    return c
}

/*Client sebarmod client */
type Client struct {
    Host string
    
    config toolkit.M
	state ClientStatus
}

/*Config return client config*/
func (c *Client) Config() toolkit.M{
    if c.config==nil {
        c.config = toolkit.M{}
    }
    
    return c.config
}

/*Status get client connection status*/
func (c *Client) Status() ClientStatus{
    return c.state
}

/*Connect connect to mod server*/
func (c *Client) Connect() error {
    if c.Host=="" {
        return errors.New("client.Connect: Host is empty")
    }
	return nil
}

/*Call call fn on server*/
func (c *Client) Call(name string, data toolkit.M) *toolkit.Result {
	ret := toolkit.NewResult()
    return ret.SetErrorTxt("mod.Client.Call: no data is being returned")
}

/*IsConnected Check if client is connected */
func (c *Client) IsConnected() bool{
    connected := c.Status()==ClientConnect
    return connected
}