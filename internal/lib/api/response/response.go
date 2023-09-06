package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const AliasLenght = 6
const (
	StatusOk    = "OK"
	StatusError = "Error"
)

func Ok() Response {
	return Response{Status: StatusOk}
}
func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}
func ValidatorError(errs validator.ValidationErrors) Response {
	var errMsg []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsg = append(errMsg, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsg = append(errMsg, fmt.Sprintf("field %s is a valid field", err.Field()))
		default:
			errMsg = append(errMsg, fmt.Sprintf("field %s is  valid ", err.Field()))

		}
	}
	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsg, ","),
	}
}
