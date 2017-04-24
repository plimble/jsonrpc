package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"sync"

	"github.com/plimble/errors"
)

type ErrorHandler func()
type MiddlewareFunc func()

type Rpc struct {
	pool          sync.Pool
	premiddleware []MiddlewareFunc
	middleware    []MiddlewareFunc
	services      *serviceMap
	Server        *http.Server
	TLSServer     *http.Server
	Listener      net.Listener
	TLSListener   net.Listener
	ErrorHandler  ErrorHandler
	Mutex         sync.RWMutex
}

func New() *Rpc {
	rpc := &Rpc{
		Server:    new(http.Server),
		TLSServer: new(http.Server),
		services:  new(serviceMap),
	}
	rpc.Server.Handler = rpc
	rpc.TLSServer.Handler = rpc

	return rpc
}

func (rpc *Rpc) Use(middleware ...MiddlewareFunc) {
	rpc.middleware = append(rpc.middleware, middleware...)
}

func (rpc *Rpc) Listen(addr string) {
	rpc.Server.Addr = addr
	rpc.Server.ListenAndServe()
}

func (rpc *Rpc) RegisterService(service interface{}, name string) {
	if err := rpc.services.register(service, name); err != nil {
		panic(err)
	}
}

func (rpc *Rpc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method != "POST" {
		rpc.writeError(w, errors.Newh(405, r.Method+" is not allowed"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if r.Header.Get("X-Rpc") == "Batch" {
		reqs := make(Requests, 0)
		if err = json.Unmarshal(body, &reqs); err != nil {
			rpc.writeError(w, errors.InternalServerError(err.Error()))
			return
		}

		ress := make(Responses, len(reqs))
		var wg sync.WaitGroup

		for i, req := range reqs {
			wg.Add(1)
			go func(idx int, req *Request) {
				ress[idx] = rpc.runRequest(req, w, r)
				wg.Done()
			}(i, req)
		}
		wg.Wait()

		rpc.write(w, ress)
		return
	}
	req := new(Request)
	if err = json.Unmarshal(body, req); err != nil {
		rpc.writeError(w, errors.InternalServerError(err.Error()))
		return
	}
	res := rpc.runRequest(req, w, r)
	rpc.write(w, res)
}

func (rpc *Rpc) runRequest(req *Request, w http.ResponseWriter, r *http.Request) *Response {
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

func (rpc *Rpc) write(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	resJson, _ := json.Marshal(res)
	w.Write(resJson)
}

func (rpc *Rpc) writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	status, err := errors.ErrorStatus(err)
	w.WriteHeader(status)
	res := &Response{
		Error: &ResponseError{
			Code:    status,
			Message: err.Error(),
		},
	}
	resJson, _ := json.Marshal(res)
	w.Write(resJson)
	// if s.afterFunc != nil {
	// 	s.afterFunc(&RequestInfo{
	// 		Error:      fmt.Errorf(msg),
	// 		StatusCode: status,
	// 	})
	// }
}
