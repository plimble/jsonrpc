package main

import (
	"context"

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
	j.RegisterService(new(Adder), "adder")

	j.Listen(":3000")
}
