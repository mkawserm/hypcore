package core

import (
	"strings"
)

type ServiceRegistry struct {
	ServiceMap map[string]ServiceInterface
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		ServiceMap: make(map[string]ServiceInterface),
	}
}

func (sr *ServiceRegistry) RegisterService(si ServiceInterface) {
	sr.ServiceMap[strings.ToLower(si.ServiceName())] = si
}

func (sr *ServiceRegistry) Get(name string) (ServiceInterface, bool) {
	si, ok := sr.ServiceMap[strings.ToLower(name)]

	return si, ok
}

func (sr *ServiceRegistry) SetServiceSettings(name string, data []byte) bool {
	si, ok := sr.ServiceMap[strings.ToLower(name)]
	if ok {
		si.SetServiceSettings(data)
		return true
	} else {
		return false
	}
}
