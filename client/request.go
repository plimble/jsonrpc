package client

import (
	"encoding/json"

	"github.com/renstrom/shortuuid"
)

type Request struct {
	Id     string          `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

func NewRequest(method string, params Params) *Request {
	pb, _ := json.Marshal(params)
	return &Request{
		Id:     shortuuid.New(),
		Method: method,
		Params: pb,
	}
}

type Params map[string]interface{}

type Requests []*Request
