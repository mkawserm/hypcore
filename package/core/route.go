package core

type Route struct {
	Pattern           string
	HttpHandlerObject ServeHTTPInterface
}
