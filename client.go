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

/*Close close connection*/
func (c *Client) Close(){
    if c.rpcclient!=nil {
        c.rpcclient.Close()
    }
}

/*CallResult call fn on server and return its value as into toolkit.Result*/
func (c *Client) CallResult(methodName string, in toolkit.M) *toolkit.Result {
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

/*Call call a function and return its result into a pointer object*/
func (c *Client) Call(methodname string, in toolkit.M, out interface{}) error{
    r := c.CallResult(methodname, in)
    if r.Status!=toolkit.Status_OK{
        return errors.New("client.CallTo: " + r.Message)
    }
    
    if out==nil{
        return nil
    }
    
    var e error
    if !r.IsEncoded() {
        e = r.Cast(out, r.EncoderID)
        if e!=nil {
            return errors.New("client.CallTo: Cast Fail " + e.Error())
        }
        return nil
    }
    
    e = r.GetFromBytes(out)
    if e!=nil{
        return errors.New("client.CallTo: Decode bytes fail " + e.Error())
    }
    
    return nil
}

/*IsConnected Check if client is connected */
func (c *Client) IsConnected() bool{
    connected := c.Status()==ClientConnect
    return connected
}