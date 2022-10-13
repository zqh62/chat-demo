package api

import (
	"chat/serializer"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
)

// ErrorResponse 返回错误信息 ErrorResponse
func ErrorResponse(err error) serializer.Response {
	if _, ok := err.(validator.ValidationErrors); ok {
		return serializer.Response{
			Status: 400,
			Msg:    "Wrong parameters.",
			Error:  fmt.Sprint(err),
		}
	}
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return serializer.Response{
			Status: 400,
			Msg:    "JSON type doesn't match.",
			Error:  fmt.Sprint(err),
		}
	}

	return serializer.Response{
		Status: 400,
		Msg:    "Wrong parameters.",
		Error:  fmt.Sprint(err),
	}
}
