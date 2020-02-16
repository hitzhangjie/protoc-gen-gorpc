package gorpc

import (
	"github.com/hitzhangjie/protoc-gen-gorpc/gorpc/gotpl"
	"github.com/hitzhangjie/protoc-gen-gorpc/gorpc/gotpl/rpc"
)

var GoRPCTemplates = map[string]string{
	"main.go":             gotpl.MainGo,
	"service_rpc.go":      gotpl.ServiceRPCGo,
	"service_rpc_test.go": gotpl.ServiceRPCTestGo,
	"go.mod":              gotpl.GoModGo,
	"rpc/go.mod":          rpc.GoModGo,
	"rpc/gorpc.go":        rpc.GoRPCGo,
}
