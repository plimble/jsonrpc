package client

import (
	"encoding/json"
)

type Response struct {
	Id     string          `json:"id,omitempty"`
	Error  *ResponseError  `json:"error,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
}

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *Response) UnmarshalResult(v interface{}) error {
	return json.Unmarshal(c.Result, v)
}

type Responses []*Response
