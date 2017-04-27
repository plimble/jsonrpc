# jsonrpc
Golang JSONRPC

## Install

#### Install Server

```
go get -u github.com/plimble/jsonrpc
```

#### Install client

```
go get -u github.com/plimble/jsonrpc/client
```

#### Install Both

```
go get -u github.com/plimble/jsonrpc/...
```

## Example Server

```go
package main

import (
	"context"

	"github.com/labstack/echo"
	"github.com/plimble/jsonrpc"
)

type Adder struct{}

type AddReq struct {
	A, B int
}

type AddRes struct {
	Val int
}

func (a *Adder) Add(ctx context.Context, req *AddReq) (*AddRes, error) {
	val := req.A + req.B
	return &AddRes{
		Val: val,
	}, nil
}

func (a *Adder) Multiply(ctx context.Context, req *AddReq) (*AddRes, error) {
	val := req.A * req.B
	return &AddRes{
		Val: val,
	}, nil
}

func main() {
	j := jsonrpc.New()
	j.Register(new(Adder), "adder")
	e := echo.New()
	e.POST("/", j.Handle)
	e.POST("/batch", j.HandleBatch)

	e.Start(":3000")
}
```

## Example Client

```go
package main

import (
	"fmt"

	"github.com/plimble/jsonrpc/client"
)

func main() {
	c := client.New("http://localhost:3000", "/", "/batch")
	res, err := c.Request("adder.Add", client.Params{"a": 1, "b": 4})
	if err != nil {
		panic(err)
	}
	if res.Error != nil {
		panic(res.Error)
	}

	result := make(map[string]interface{})
	res.UnmarshalResult(&result)
	fmt.Println("Result:", result)

	ress, err := c.Requests(&client.Requests{
		client.NewRequest("adder.Add", client.Params{"a": 1, "b": 2}),
		client.NewRequest("adder.Add", client.Params{"a": 2, "b": 2}),
		client.NewRequest("adder.Add", client.Params{"a": 3, "b": 2}),
		client.NewRequest("adder.Multiply", client.Params{"a": 4, "b": 2}),
	})
	if err != nil {
		panic(err)
	}

	for _, res := range ress {
		result := make(map[string]interface{})
		res.UnmarshalResult(&result)
		fmt.Println("Batch Result:", result)
	}
}
```

