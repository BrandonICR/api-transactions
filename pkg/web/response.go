package web

import "fmt"

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func NewResponse(code int, message string, data interface{}, err string) Response {
	if code < 400 {
		return Response{fmt.Sprint(code), message, data, ""}
	}
	return Response{fmt.Sprint(code), message, nil, err}
}
