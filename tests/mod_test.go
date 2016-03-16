package tests

import (
    "github.com/eaciit/toolkit"
    "github.com/eaciit/sebarmod"
    "time"
    "errors"
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

func TestCall(t *testing.T){
    client = sebarmod.NewClient("localhost:5000", nil)
    e = client.Connect()
    check("CallConnect", e, t)

var result *toolkit.Result
    result = client.Call("hello",toolkit.M{}.Set("name","Arief Darmawan"))
    if result.Status!=toolkit.Status_OK{
        check("Call", errors.New(result.Message), t)
    }
    
    var returned *ModObj
    if e=result.GetFromBytes(&returned); e!=nil {
        check("DecodeReturn", e, t)
    }
    toolkit.Println("Value returned:\n", toolkit.JsonStringIndent(returned,"\t"))
}

func TestClose(t *testing.T){
    skipIfNil(t)
    e = svr.Stop()
    check("Stop", e, t)
}