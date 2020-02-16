package gotpl

var GoModGo=`{{- $pkgName := .PackageName -}}
{{- $svrName := (index .Services 0).Name -}}

{{- if eq .GoMod "" -}}
module gorpc.app.{{$svrName}}
{{- else -}}
module {{.GoMod}}
{{- end }}

go 1.12

replace {{$pkgName}} => ./{{$pkgName}}`

