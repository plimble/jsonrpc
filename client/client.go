package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/plimble/errors"
	"github.com/renstrom/shortuuid"
)

type Client struct {
	addr string
}

func New(addr string) *Client {
	return &Client{addr}
}

func (c *Client) Request(method string, params Params) (*Response, error) {
	pb, _ := json.Marshal(params)
	req := &Request{
		Id:     shortuuid.New(),
		Method: method,
		Params: pb,
	}

	reqJson, _ := json.Marshal(req)
	httpreq, _ := http.NewRequest("POST", c.addr, bytes.NewBuffer(reqJson))
	httpreq.Header.Set("Content-Type", "application/json; charset=utf-8")
	res, err := http.DefaultClient.Do(httpreq)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	defer res.Body.Close()

	cres := new(Response)
	err = json.Unmarshal(body, cres)

	return cres, err
}

func (c *Client) Requests(reqs *Requests) (Responses, error) {
	reqsJson, _ := json.Marshal(reqs)
	httpreq, _ := http.NewRequest("POST", c.addr, bytes.NewBuffer(reqsJson))
	httpreq.Header.Set("Content-Type", "application/json; charset=utf-8")
	httpreq.Header.Set("X-Rpc", "Batch")
	res, err := http.DefaultClient.Do(httpreq)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	defer res.Body.Close()

	cres := make(Responses, 0)
	err = json.Unmarshal(body, &cres)

	return cres, err
}
