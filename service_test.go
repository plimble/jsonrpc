package jsonrpc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"golang.org/x/net/context"
)

type Adder struct{}

type AddReq struct {
	A, B int
}

type AddRes struct {
	Val int
}

func (a *Adder) Add(ctx context.Context, req *AddReq) (*AddRes, error) {
	fmt.Println(ctx.Value("user"))
	val := req.A + req.B
	return &AddRes{
		Val: val,
	}, nil
}

func TestMap(t *testing.T) {
	a := new(Adder)

	sm := new(serviceMap)
	err := sm.register(a, "adder")
	if err != nil {
		t.Error(err)
	}

	ss, ms, _ := sm.get("adder.Add")

	data, _ := json.Marshal(&AddReq{1, 2})

	req := reflect.New(ms.argsType)
	json.Unmarshal(data, req.Interface())

	vv := context.Background()
	aa := context.WithValue(vv, "user", "xxxx")
	ctx := reflect.ValueOf(aa)

	var retVals []reflect.Value
	retVals = ms.method.Func.Call([]reflect.Value{
		ss.rcvr,
		ctx,
		req,
	})

	result, _ := json.Marshal(retVals[0].Interface())

	fmt.Println(string(result))

}
