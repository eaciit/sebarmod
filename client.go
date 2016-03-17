package sebarmod

import (
	"github.com/eaciit/toolkit"
    "errors"
    "net/rpc"
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
    
    rpcclient *rpc.Client
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
    
    client, econnect := rpc.Dial("tcp", c.Host)
    if econnect!=nil {
        return errors.New("client.Connect: " + econnect.Error())
    }
    c.rpcclient = client
	return nil
}

/*Call call fn on server*/
func (c *Client) Call(methodName string, in toolkit.M) *toolkit.Result {
	if c.rpcclient == nil {
		return toolkit.NewResult().SetErrorTxt(toolkit.Sprintf("Unable to call %s.%s, no connection handshake", c.Host, methodName))
	}
	if in == nil {
		in = toolkit.M{}
	}
	out := toolkit.NewResult()
	in["method"] = methodName
	e := c.rpcclient.Call("RPC.Do", in, out)
	if e != nil {
		return out.SetErrorTxt(c.Host + "." + methodName + " Fail: " + e.Error())
	}
	return out
}

/*IsConnected Check if client is connected */
func (c *Client) IsConnected() bool{
    connected := c.Status()==ClientConnect
    return connected
}