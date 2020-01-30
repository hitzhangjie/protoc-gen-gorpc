{{- $pkgName := .PackageName -}}
{{- $svrName := (index .Services 0).Name -}}

{{- $goPkgOption := "" -}}
{{- with .FileOptions.go_package -}}
  {{- $goPkgOption = . -}}
{{- end -}}

{{- if eq .GoMod "" -}}
module gorpc.app.{{$svrName}}
{{- else -}}
module {{.GoMod}}
{{- end }}

go 1.12

{{ $rpcdir := "" -}}
{{ if ne $goPkgOption "" -}}
{{ $rpcdir = $goPkgOption }}
{{- else -}}
{{ $rpcdir = $pkgName }}
{{- end -}}
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
{{ end }}

