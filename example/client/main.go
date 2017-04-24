package main

import (
	"fmt"

	"github.com/plimble/jsonrpc/client"
)

func main() {
	c := client.New("http://localhost:3000")
	res, err := c.Request("adder.Add", client.Params{"a": 1, "b": 4})
	if err != nil {
		panic(err)
	}
	if res.Error != nil {
		panic(res.Error)
	}

	result := make(map[string]interface{})
	fmt.Println("@@@", res)
	res.UnmarshalResult(&result)
	fmt.Println("###", result)

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
		fmt.Println("@@@", res)
		res.UnmarshalResult(&result)
		fmt.Println("###", result)
	}
}
