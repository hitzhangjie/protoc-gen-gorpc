package gotpl

var GoModGo=`{{- $pkgName := .PackageName -}}
{{- $svrName := (index .Services 0).Name -}}

{{- if eq .GoMod "" -}}
module gorpc.app.{{$svrName}}
{{- else -}}
module {{.GoMod}}
{{- end }}

go 1.12

{{ $rpcdir := $pkgName }}
replace {{$rpcdir}} => ./stub/{{$rpcdir}}

{{ range $k, $v := .Pb2ImportPath -}}

{{ if and (ne $v "") (ne $v "github.com/hitzhangjie/go-rpc/whisper") -}}
{{ $depdir := "" -}}

{{ if (contains $k "/") -}}
{{ $depdir = (trimright "/" $k) -}}
{{ end -}}

{{ if and (ne $depdir "") (ne $depdir $pkgName) -}}
replace {{$v}} => ./stub/{{$rpcdir}}/{{$depdir}}
{{ end }}

{{ if and (eq $depdir "") (ne $v $pkgName) }}
replace {{$v}} => ./stub/{{$rpcdir}}/{{$v}}
{{ end }}

{{ end }}
{{ end }}`
