package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error interface {
	Code() string
	Message() string
}

type httpError struct {
	code    string
	message string
}

func (e httpError) Error() string {
	return fmt.Sprintf("[%s] %s", e.code, e.message)
}
func (e httpError) Code() string    { return e.code }
func (e httpError) Message() string { return e.message }

func newResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	switch body.(type) {
	case []byte, string:
		fmt.Fprintf(w, "%s", body)
	default:
		json.NewEncoder(w).Encode(body)
	}
}

func newErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	body := make(map[string]interface{})
	if httpErr, ok := err.(Error); ok {
		body["code"] = httpErr.Code()
		body["message"] = httpErr.Message()
	} else {
		body["code"] = "unknown"
		body["message"] = err.Error()
	}

	newResponse(w, statusCode, body)
}
