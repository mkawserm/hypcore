package core

import (
	"net/http"
)

type ServeHTTPInterface interface {
	ServeHTTP(ctx *HContext, request *http.Request, response http.ResponseWriter)
}
