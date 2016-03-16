package sebarmod

import (
    "reflect"
    "errors"
    "github.com/eaciit/toolkit"
    "strings"
)

/*RPCFn function contract*/

/*Register register an object into RPC Server*/
func (a *Server) Register(o interface{}) error {
	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)
	if v.Kind() != reflect.Ptr {
		return errors.New("Invalid object for RPC Register")
	}
	if a.Log == nil {
		a.Log, _ = toolkit.NewLog(true, false, "", "", "")
	}
	
    objName := toolkit.TypeName(o)
	methodCount := t.NumMethod()
	for i := 0; i < methodCount; i++ {
		method := t.Method(i)
		mtype := method.Type
		methodName := strings.ToLower(method.Name)
		//fmt.Println("Evaluating " + toolkit.TypeName(o) + "." + methodName)

		//-- now check method signature
		if mtype.NumIn() == 2 && mtype.In(1).String() == "toolkit.M" {
			if mtype.NumOut() == 1 && mtype.Out(0).String() == "*toolkit.Result" {
				a.Log.Info("Registering function " + objName + "." + methodName)
				a.AddFn(methodName, v.Method(i).Interface().(func(toolkit.M) *toolkit.Result))
			}
		}
	}
	return nil
}

/*AddFn add a function to RPC server
Function should follow contract: func(in toolkit.M) *toolkit.Result
*/
func (a *Server) AddFn(methodname string, fn RPCFn) {
	var r *RPC
	if a.rpcObject == nil {
		r = new(RPC)
	} else {
		r = a.rpcObject
	}
	addFnToRPC(r, a, methodname, fn, nil)
	a.rpcObject = r
}