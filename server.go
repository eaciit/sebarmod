package sebarmod

import (
	"errors"
	"net"
	"net/rpc"
	"strings"
	"sync"

	"github.com/eaciit/toolkit"
)

type sebarFn struct {
	fn func()

	Broadcastable BroadcastTo
}

/*Server SebarMod server */
type Server struct {
	Host string
	Log  *toolkit.LogEngine

	rpcObject    *RPC
	rpcServer    *rpc.Server
	masterHost   string
	masterClient *Client
	listener     net.Listener
	fns          map[string]*sebarFn
	//nodes map[string]*Server
	clients map[string]*Client
}

/*Start start the server*/
func (s *Server) Start() error {
	everify := s.Verify()
	if everify != nil {
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

	s.AddFn("follow", func(in toolkit.M) *toolkit.Result {
		result := toolkit.NewResult()
		nodeid := in.GetString("nodeid")
		if nodeid == "" {
			return result.SetErrorTxt("nodeid should not be empty")
		}
		nodeclient := NewClient(nodeid, nil)
		econnect := nodeclient.Connect()
		if econnect != nil {
			return result.SetErrorTxt("Can not handshake with client node. " + econnect.Error())
		}
		if s.clients == nil {
			s.clients = map[string]*Client{}
		}
		s.clients[nodeid] = nodeclient
		s.Log.AddLog(toolkit.Sprintf("%s has new follower %s", s.Host, nodeid), "INFO")
      	return result
	})

	s.AddFn("unfollow", func(in toolkit.M) *toolkit.Result {
		result := toolkit.NewResult()
		nodeid := in.GetString("nodeid")
		if nodeid == "" {
			return result.SetErrorTxt("nodeid should not be empty")
		}
		if s.clients == nil {
			s.clients = map[string]*Client{}
		}
		if c, hasClient := s.clients[nodeid]; hasClient {
			c.Close()
			delete(s.clients, nodeid)
		}
		s.Log.AddLog(toolkit.Sprintf("Node %s is not following %s any longer", nodeid, s.Host), "INFO")
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
func (s *Server) Stop() error {
	if s.clients != nil {
		for _, c := range s.clients {
			rstop := c.CallResult("stop", nil)
			if rstop.Status != toolkit.Status_NOK {
				return errors.New(rstop.Message)
			}
		}
	}

	if s.masterClient != nil {
		s.masterClient.Close()
	}
	return nil
}

/*Verify verify server*/
func (s *Server) Verify() error {
	return nil
}

func (s *Server) initFn() {
	if s.fns == nil {
		s.fns = map[string]*sebarFn{}
	}
}

/*SetFn set function*/
func (s *Server) SetFn(name string, fn func(), config toolkit.M) {
	if fn == nil {
		return
	}
    
    if config==nil{
        config=toolkit.M{}
    }
    
    broadcast := config.Get("broadcast", 0).(BroadcastTo)
    
  	s.initFn()
	s.fns[name] = &sebarFn{
		fn: fn,
        Broadcastable: broadcast,
	}
}

func (s *Server) initClients() {
	if s.clients == nil {
		s.clients = map[string]*Client{}
	}
}

func (s *Server) client(id string) *Client {
	s.initClients()
	c, has := s.clients[id]
	if !has {
		return nil
	}
	if !c.IsConnected() {
		econnect := c.Connect()
		if econnect != nil {
			return nil
		}
	}
	return c
}

/*BroadcastTo To determine how a function will be run*/
type BroadcastTo int

const (
	/*BroadcastAll = Broadcast to all nodes*/
	BroadcastAll BroadcastTo = 0
	/*BroadcastOnly = BroadcastOnly to specified nodes*/
	BroadcastOnly = 1
	/*BroadcastExcept  = Broadcast to all nodes except to specified ones*/
	BroadcastExcept = 2
)

/*Broadcast broadcast call to all node*/
func (s *Server) broadcast(broadcastto BroadcastTo, nodeids []string, name string, config toolkit.M) (success, fail map[string]*toolkit.Result) {
	success = map[string]*toolkit.Result{}
	fail = map[string]*toolkit.Result{}
	if s.clients == nil {
		return
	}
	wg := new(sync.WaitGroup)
	for id, c := range s.clients {
		wg.Add(1)
		go func(id string, c *Client, wg *sync.WaitGroup) {
			defer wg.Done()
			result := c.CallResult(name, config)
			if result.Status == toolkit.Status_OK {
				success[id] = result
			} else {
				fail[id] = result
			}
		}(id, c, wg)
	}
	wg.Wait()
	return
}

/*SetMaster set master server*/
func (s *Server) SetMaster(host string, config toolkit.M) error {
	if s.masterClient != nil {
        s.masterClient.Call("unfollow", toolkit.M{}.Set("nodeid", s.masterClient.Host), nil)
		s.masterClient.Close()
	}
    
    s.masterHost = host
	if host==""{
        return nil    
    }
    
    masterClient := NewClient(s.masterHost, config)
	e := masterClient.Connect()
	if e != nil {
		return errors.New("Server.SetMaster: " + e.Error())
	}

	e = masterClient.Call("follow", toolkit.M{}.Set("nodeid", s.Host), nil)
	if e != nil {
		return errors.New("Server.SetMaster.Follow: " + e.Error())
	}
    s.Log.AddLog(toolkit.Sprintf("Server %s is now following %s", s.Host, host), "INFO")  	
	return nil
}

/*Master get master server host*/
func (s *Server) Master() string {
	return s.masterHost
}
