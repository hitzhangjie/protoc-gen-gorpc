package gotpl

var MainGo = `{{- $pkgName := .PackageName -}}
package main

import (
	gorpc "github.com/hitzhangjie/go-rpc"

    pb "{{$pkgName}}"
)

{{range $index, $service := .Services}}
{{- $svrName := $service.Name | camelcase | untitle -}}
type {{$svrName}}ServiceImpl struct {}
{{end}}

func main() {

	s := gorpc.NewServer()

    {{range $index, $service := .Services}}
    {{- $svrNameCamelCase := $service.Name | camelcase -}}
	pb.Register{{$svrNameCamelCase}}Service(s, &{{$svrNameCamelCase|untitle}}ServiceImpl{})
	{{end}}
	s.Serve()
}`
