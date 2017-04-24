package jsonrpc

type Response struct {
	Id     string         `json:"id,omitempty"`
	Error  *ResponseError `json:"error,omitempty"`
	Result interface{}    `json:"result,omitempty"`
}

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r *Response) SetErr(code int, message string) {
	r.Error = &ResponseError{
		Code:    code,
		Message: message,
	}
}

type Responses []*Response
