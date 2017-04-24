package jsonrpc

import (
	"encoding/json"

	"github.com/renstrom/shortuuid"
)

type Request struct {
	Id     string          `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type Params map[string]interface{}

func NewRequest(method string, params Params) Request {
	pb, _ := json.Marshal(params)
	return Request{
		Id:     shortuuid.New(),
		Method: method,
		Params: pb,
	}
}

type Requests []*Request
