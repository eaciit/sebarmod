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
    Value int
}

func newModObj(x int) *ModObj{
    return &ModObj{time.Now(), x}
}

type ModApp struct{    
}

func (m *ModApp) Hello(name string) string{
    return toolkit.Sprintf("Hello %s", name)
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

func TestClose(t *testing.T){
    skipIfNil(t)
    e = svr.Stop()
    check("Stop", e, t)
}