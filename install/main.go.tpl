{{- $pkgName := .PackageName -}}
{{- $goPkgOption := "" -}}
{{- with .FileOptions.go_package -}}
  {{- $goPkgOption = . -}}
{{- end -}}
package main

import (
	gorpc "github.com/hitzhangjie/go-rpc"

    {{ if ne $goPkgOption "" -}}
   	pb "{{$goPkgOption}}"
    {{- else -}}
    pb "{{$pkgName}}"
	{{- end }}
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
}
