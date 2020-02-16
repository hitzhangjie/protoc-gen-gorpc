package gotpl

var ServiceRPCGo = `{{- $pkgName := .PackageName -}}
{{- $imports := .Imports -}}
package main

import (
	"context"

	"{{$pkgName -}}"

{{ range $imports }}
{{ . }}
{{ end -}}
)

{{ $service := (index .Services .ServiceIndex) -}}
{{- $serviceName := $service.Name | camelcase | untitle -}}

{{ range $index, $method := $service.RPC }}
{{- $rpcName := $method.Name | camelcase -}}
{{- $rpcReqType := $method.RequestType -}}
{{- $rpcRspType := $method.ResponseType -}}

func (s *{{$serviceName}}ServiceImpl) {{$rpcName}}(ctx context.Context, req *{{$rpcReqType}}, rsp *{{$rpcRspType}}) (err error) {
	// implement business logic here ...
	// ...

	return
}

{{end}}`
