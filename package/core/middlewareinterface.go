package core

import (
	"net/http"
)

type MiddlewareInterface interface {
	// ServeHTTP and instruct to proceed next or not
	ServeHTTP(context *HContext, request *http.Request, response http.ResponseWriter) bool
}
