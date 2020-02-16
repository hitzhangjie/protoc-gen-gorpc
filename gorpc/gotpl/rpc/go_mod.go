package rpc

var GoModGo = `{{- $pkgName := .PackageName -}}
module {{$pkgName}}

go 1.12`
