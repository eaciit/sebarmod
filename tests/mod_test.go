package tests

import (
    "github.com/eaciit/toolkit"
    "github.com/eaciit/sebarmod"
    "time"
    "testing"
)

var (
    e error
    svr *sebarmod.Server
)

type ModObj struct{
    Created time.Time
    Name string
    Value int
}

func newModObj(name string) *ModObj{
    return &ModObj{time.Now(), name, len(name)}
}

type ModApp struct{    
}

func (m *ModApp) Hello(in toolkit.M) *toolkit.Result{
    r := toolkit.NewResult()
    data := newModObj(in.GetString("name"))
    return r.SetBytes(data, "")
}

func skipIfNil(t *testing.T){
    if svr==nil {
        t.Skip()
    }
}

func check(pre string, e error, t *testing.T){
    if e!=nil {
        t.Fatal(pre, e.Error())
    }
}

func TestServer(t *testing.T){
    svr = sebarmod.NewServer("localhost:5000")
    e = svr.Register(new(ModApp))
    check("Register ModApp", e, t)
    e = svr.Verify()
    check("Start Server", e, t)
}

func TestStart(t *testing.T){
    skipIfNil(t)
    e = svr.Start()
    check("Start",e,t)
}

var client *sebarmod.Client

func TestPing(t *testing.T){
    client = sebarmod.NewClient("localhost:5000", nil)
    e = client.Connect()
    check("CallConnect", e, t)

    returned := ""
    e = client.Call("ping",nil,&returned)
    check("Call", e, t)
    toolkit.Println("Value returned:\n", toolkit.JsonStringIndent(returned,"\t"))
}

func TestCall(t *testing.T){
    returned := new(ModObj)
    e = client.Call("hello",toolkit.M{}.Set("name","Arief Darmawan Soebani"), returned)
    check("Call", e, t)
    toolkit.Println("Value returned:\n", toolkit.JsonStringIndent(returned,"\t"))
}

var nodes []*sebarmod.Server
func TestNode(t *testing.T){
    for i:=0;i<5;i++{
        port := 5001+i
        host := toolkit.Sprintf("localhost:%d",port)
        snode := sebarmod.NewServer(host)
        snode.Register(new(ModApp))
        snode.Start()
        efollow := snode.SetMaster("localhost:5000", nil)
        if efollow!=nil {
            t.Fatalf("Fail to follow master on node %s: %s", host, efollow.Error())
        }
        nodes = append(nodes, snode)
    }
}

func TestClose(t *testing.T){
    skipIfNil(t)
    e = svr.Stop()
    check("Stop", e, t)
    
    /*
    for _, node := range nodes{
        node.Stop()
    }
    */
    
    client.Close()
}