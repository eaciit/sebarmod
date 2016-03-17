package sebarmod

import (
    "github.com/eaciit/toolkit"
    "net"
    "net/rpc"
    "errors"
    "strings"
)

type sebarFn struct {
	fn func()
    
    Broadcastable bool
}

/*Server SebarMod server */
type Server struct {
    Host string
    Log *toolkit.LogEngine
    
    rpcObject *RPC
	rpcServer *rpc.Server
    listener net.Listener
	fns   map[string]*sebarFn
	nodes map[string]*Server
    clients map[string]*Client
}

/*Start start the server*/
func (s *Server) Start() error {
    everify := s.Verify()
    if everify!=nil {
        return everify
    }
    
    if s.Log == nil {
		le, e := toolkit.NewLog(true, false, "", "", "")
		if e == nil {
			s.Log = le
		} else {
			return errors.New("Unable to setup log")
		}
	}

	if s.rpcObject == nil {
		s.rpcObject = new(RPC)
	}
	
	s.AddFn("ping", func(in toolkit.M) *toolkit.Result {
		result := toolkit.NewResult()
		result.Data = "Application Server powered by SebarMod"
		return result
	})

	s.Log.Info("Starting server " + s.Host + ". Registered functions are: " + strings.Join(func() []string {
		ret := []string{}
		for k := range s.rpcObject.Fns {
			ret = append(ret, k)
		}
		return ret
	}(), ", "))

	s.rpcServer = rpc.NewServer()
	s.rpcServer.Register(s.rpcObject)
	l, e := net.Listen("tcp", toolkit.Sprintf("%s", s.Host))
	if e != nil {
		return e
	}

	s.listener = l
	go func() {
		s.rpcServer.Accept(l)
	}()
	return nil
}

/*Stop stop the server*/
func (s *Server) Stop()error{
    for id := range s.nodes{
        if c:=s.client(id); c!=nil {
            rstop := c.CallResult("stop", nil)
            if rstop.Status!=toolkit.Status_NOK{
                return errors.New(rstop.Message)
            }
        }
    }
    return nil
}

/*Verify verify server*/
func (s *Server) Verify()error {
    return nil
}

func (s *Server) initFn() {
	if s.fns == nil {
		s.fns = map[string]*sebarFn{}
	}
}

/*SetFn set function*/
func (s *Server) SetFn(name string, fn func()) {
	if fn == nil {
		return
	}

	s.initFn()
	s.fns[name] = &sebarFn{
		fn: fn,
	}
}

func (s *Server) initNodes() {
	if s.nodes == nil {
		s.nodes = map[string]*Server{}
	}
}

/*AddNode add a server node*/
func (s *Server) AddNode(nodeid string, nodeservers ...*Server) {
	s.initNodes()
	for _, nodeserver := range nodeservers {
		s.nodes[nodeid] = nodeserver
	}
}

/*RemoveNode remove server node*/
func (s *Server) RemoveNode(nodeids ...string) {
	if s.nodes == nil {
		return
	}
	for _, nodeid := range nodeids {
		delete(s.nodes, nodeid)
	}
}

/*Node Get node from server nodes */
func (s *Server) Node(id string) *Server{
    if s.nodes==nil {
        return nil
    }
    return s.nodes[id]
}

func (s *Server) client(id string) *Client{
    if s.clients==nil {
        s.clients=map[string]*Client{}
    }
    n := s.Node(id)
    if n==nil {
       return nil 
    }
    c, has := s.clients[id]
    if !has{
        return nil
    }
    if !c.IsConnected(){
        econnect := c.Connect()
        if econnect!=nil {
            return nil
        }
    }
    return c
}

type BroadcastTo int
const (
    BroadcastAll BroadcastTo = 0
    BroadcastOnly = 1
    BroadcastExcept = 2
)

/*Broadcast broadcast call to all node*/
func (s *Server) broadcast(broadcastto BroadcastTo, nodeids []string, name string){
    if s.nodes==nil {
        return 
    }
    for id, _ := range s.nodes{
        c := s.client(id)
        if c!=nil {
            result := c.CallResult(name, nil)
            if result.Status!=toolkit.Status_OK{
                
            }
         }
    }
}