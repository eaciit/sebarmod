package sebarmod

import (
    "github.com/eaciit/toolkit"
)

/*Client sebarmod client */
type Client struct{
}

/*Connect connect to mod server*/
func (c *Client) Connect(server string) error{
    return nil
}

/*Call call fn on server*/
func (c *Client) Call(name string, data toolkit.M, output interface{})error{
    return nil
}