package jsonrpc

import (
	"encoding/json"
)

type Request struct {
	Id     string          `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type Params map[string]interface{}

type Requests []*Request
