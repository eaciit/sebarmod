package sebarmod

type SebarFn func()

/*Server SebarMod server */
type Server struct{
    fns map[string]SebarFn
    nodes map[string]*Server
}

/*Start start the server*/
func (s *Server) Start() error{
    return nil
}

func (s *Server) initFn(){
    if s.fns==nil {
        s.fns=map[string]SebarFn{}
    }
}

/*SetFn set function*/
func (s *Server) SetFn(name string, fn SebarFn){
    if fn==nil {
        return
    }
    
    s.initFn()
    s.fns[name]=fn
}

func (s *Server) initNodes(){
    if s.nodes==nil {
        s.nodes=map[string]*Server{}
    }
}

/*AddNode add a server node*/
func (s *Server) AddNode(nodeid string, nodeservers ...*Server){
    s.initNodes()
    for _, nodeserver := range nodeservers{
        s.nodes[nodeid]=nodeserver
    }
}

/*RemoveNode remove server node*/
func (s *Server) RemoveNode(nodeids ...string){
    if s.nodes==nil {
        return
    }
    for _, nodeid := range nodeids{
        delete(s.nodes, nodeid)
    }
}