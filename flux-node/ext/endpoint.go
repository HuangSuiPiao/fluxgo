package ext

import (
	"github.com/bytepowered/flux/flux-node"
	"strings"
	"sync"
)

var (
	endpoints = new(sync.Map)
)

func MakeEndpointKey(method, pattern string) string {
	return strings.ToUpper(method) + "#" + pattern
}

func RegisterEndpoint(key string, endpoint *flux.Endpoint) *flux.MVCEndpoint {
	mvce := flux.NewMVCEndpoint(endpoint)
	endpoints.Store(key, mvce)
	return mvce
}

func EndpointByKey(key string) (*flux.MVCEndpoint, bool) {
	ep, ok := endpoints.Load(key)
	if ok {
		return ep.(*flux.MVCEndpoint), true
	}
	return nil, false
}

func Endpoints() map[string]*flux.MVCEndpoint {
	out := make(map[string]*flux.MVCEndpoint, 32)
	endpoints.Range(func(key, value interface{}) bool {
		out[key.(string)] = value.(*flux.MVCEndpoint)
		return true
	})
	return out
}
