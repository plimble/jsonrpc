package jsonrpc

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"

	"github.com/labstack/echo"
	"github.com/plimble/errors"
)

type JsonRpc struct {
	services *serviceMap
}

func New() *JsonRpc {
	return &JsonRpc{
		services: new(serviceMap),
	}
}

func (rpc *JsonRpc) Register(service interface{}, name string) {
	if err := rpc.services.register(service, name); err != nil {
		panic(err)
	}
}

func (rpc *JsonRpc) Handle(c echo.Context) error {
	var err error

	if err = rpc.checkMethod(c); err != nil {
		return err
	}

	req := new(Request)
	if err = c.Bind(req); err != nil {
		return errors.InternalServerError(err.Error())
	}
	res := rpc.runRequest(req, c)
	return c.JSON(200, res)
}

func (rpc *JsonRpc) HandleBatch(c echo.Context) error {
	var err error

	if err = rpc.checkMethod(c); err != nil {
		return err
	}

	reqs := make(Requests, 0)
	if err = c.Bind(&reqs); err != nil {
		return errors.InternalServerError(err.Error())
	}

	ress := make(Responses, len(reqs))
	var wg sync.WaitGroup

	for i, req := range reqs {
		wg.Add(1)
		go func(idx int, req *Request) {
			ress[idx] = rpc.runRequest(req, c)
			wg.Done()
		}(i, req)
	}
	wg.Wait()

	return c.JSON(200, ress)
}

func (rpc *JsonRpc) checkMethod(c echo.Context) error {
	if method := c.Request().Method; method != "POST" {
		return errors.Newh(405, method+" is not allowed")
	}

	return nil
}

func (rpc *JsonRpc) runRequest(req *Request, c echo.Context) *Response {
	var err error
	result := &Response{}
	res := &Response{
		Id: req.Id,
	}
	ss, ms, err := rpc.services.get(req.Method)
	if err != nil {
		result.SetErr(400, err.Error())
		return res
	}

	args := reflect.New(ms.argsType)
	if err = json.Unmarshal(req.Params, args.Interface()); err != nil {
		res.SetErr(500, err.Error())
		return res
	}

	ctx := context.Background()
	ctxVO := reflect.ValueOf(ctx)

	var retVals []reflect.Value
	retVals = ms.method.Func.Call([]reflect.Value{
		ss.rcvr,
		ctxVO,
		args,
	})

	errInter := retVals[1].Interface()
	if errInter != nil {
		status, err := errors.ErrorStatus(errInter.(error))
		res.SetErr(status, err.Error())
		return res
	}

	res.Result = retVals[0].Interface()

	return res
}
