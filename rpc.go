package sebarmod

import (
	"github.com/eaciit/toolkit"
	//"time"
	"errors"
	"strings"
)

type RPCFn func(toolkit.M) *toolkit.Result
type RPCFns map[string]*RPCFnInfo

type RPCFnInfo struct {
	//AuthRequired bool
	//AuthType     string
    Fn           RPCFn
    
    config toolkit.M
}

/*Config return function config*/
func (fi *RPCFnInfo) Config() toolkit.M{
    if fi.config==nil {
        fi.config = toolkit.M{}
    }
    return fi.config
}

//type ReturnedBytes []byte

type RPC struct {
	Fns               RPCFns
	Server            *Server
	MarshallingMethod string
}

var _marshallingMethod string

func MarshallingMethod() string {
	if _marshallingMethod == "" {
		_marshallingMethod = "gob"
	} else {
		_marshallingMethod = strings.ToLower(_marshallingMethod)
	}
	return _marshallingMethod
}

func SetMarshallingMethod(m string) {
	_marshallingMethod = m
}

func (r *RPC) Do(in toolkit.M, out *toolkit.Result) error {
	if r.Fns == nil {
		r.Fns = map[string]*RPCFnInfo{}
	}

	//in.Set("RPC", r)
	method := in.GetString("method")
	if method == "" {
		return errors.New("Method is empty")
	}
	fninfo, fnExist := r.Fns[method]
	if !fnExist {
		return errors.New("Method " + method + " is not exist")
		
	}
	res := fninfo.Fn(in)
	if res.Status != toolkit.Status_OK {
		return errors.New("RPC Call error: " + res.Message)
	}
	//*out = toolkit.ToBytes(res.Data, MarshallingMethod())
	*out = *res
	return nil
}

/*AddFn add a function to server */
func addFnToRPC(r *RPC, svr *Server, k string, fn RPCFn, config toolkit.M) {
	//func (r *RPC) AddFn(k string, fn RPCFn) {
	//if r.Server == nil {
	svr.Log.Info("Register " + svr.Host + "/" + k)
	r.Server = svr
	//}
	if r.Fns == nil {
		r.Fns = map[string]*RPCFnInfo{}
	}
	r.Fns[k] = &RPCFnInfo{
		//AuthRequired: needValidation,
		//AuthType:     authType,
        config: config,
		Fn:           fn,
	}
}
