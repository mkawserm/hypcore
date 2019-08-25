package core

import (
	"net/http"
)

type ServeHTTPInterface interface {
	ServeHTTP(context *HContext, request *http.Request, response http.ResponseWriter)
}
