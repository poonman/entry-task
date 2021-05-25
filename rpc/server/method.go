package server

import "context"

type MethodHandler func(srv interface{}, ctx context.Context, dec func(interface{}) error ) (rsp interface{}, err error)

type MethodDesc struct {
	Name string
	Handler MethodHandler
}


