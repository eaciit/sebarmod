package tests

import (
    "github.com/eaciit/sebarmod"
    "testing"
)

var (
    e error
    svr *sebarmod.Server
)

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
    e = svr.Verify()
    if e!=nil {
        t.Fatal("Error start server:", e.Error())
    }
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