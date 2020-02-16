package gotpl

import "github.com/hitzhangjie/protoc-gen-gorpc/gorpc/gotpl/rpc"

var GoRPCTemplates = map[string]string{
	"main.go":             MainGo,
	"service_rpc.go":      ServiceRPCGo,
	"service_rpc_test.go": ServiceRPCTestGo,
	"go.mod":              GoModGo,
	"rpc/go.mod":          rpc.GoModGo,
	"rpc/gorpc.go":        rpc.GoRPCGo,
}
