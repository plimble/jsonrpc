package client

import (
	"fmt"
	"testing"
)

func TestRequest(t *testing.T) {
	c := New("http://localhost:3000")
	res, err := c.Request("adder.Add", Params{"a": 1, "b": 4})
	if err != nil {
		t.Error(err)
	}
	if res.Error != nil {
		t.Error(res.Error)
	}

	result := make(map[string]interface{})
	fmt.Println("@@@", res)
	res.UnmarshalResult(&result)
	fmt.Println("###", result)
}

func TestRequests(t *testing.T) {
	c := New("http://localhost:3000")
	ress, err := c.Requests(&Requests{
		NewRequest("adder.Add", Params{"a": 1, "b": 2}),
		NewRequest("adder.Add", Params{"a": 2, "b": 2}),
		NewRequest("adder.Add", Params{"a": 3, "b": 2}),
		NewRequest("adder.Multiply", Params{"a": 4, "b": 2}),
	})
	if err != nil {
		t.Error(err)
	}

	for _, res := range ress {
		result := make(map[string]interface{})
		fmt.Println("@@@", res)
		res.UnmarshalResult(&result)
		fmt.Println("###", result)
	}
}
