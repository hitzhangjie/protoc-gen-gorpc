// Code generated by go-rpc-cmdline/protoc-gen-gorpc. DO NOT EDIT.
// source: {{.ProtoFile}}

{{ $pkgName := .PackageName -}}
{{- $goPkgOption := "" -}}
{{- with .FileOptions.go_package -}}
  {{- $goPkgOption = . -}}
{{- end -}}
{{- if ne $goPkgOption "" -}}
package {{ (splitList "/" $goPkgOption)|last|gopkg}}
{{- else -}}
package {{ $pkgName|gopkg }}
{{- end }}

import (
	"context"

   _ "github.com/hitzhangjie/go-rpc/"
   _ "github.com/hitzhangjie/go-rpc/http"

   	{{ if and (ne .Protocol "whisper") (ne .Protocol "http") }}
   	_ "github.com/hitzhangjie/go-rpc-codec/{{.Protocol}}"
   	{{ end }}

    "github.com/hitzhangjie/go-rpc/server"
    "github.com/hitzhangjie/go-rpc/client"
    "github.com/hitzhangjie/go-rpc/codec"

    {{ range .Imports }}
    {{ if and (ne $pkgName .) (ne $goPkgOption .) }}
    "{{ . }}"
    {{ end }}
    {{ end }}
)

/* ************************************ Service Definition ************************************ */

{{ range $service := .Services }}
{{- $svrName := $service.Name -}}
{{- $svrNameCamelCase := $service.Name|camelcase -}}
// {{$svrNameCamelCase}}Service defines service
type {{$svrNameCamelCase}}Service interface {

	{{ range $index, $method := $service.RPC }}
	{{- $rpcName := $method.Name | camelcase -}}
    {{- $rpcReqType := $method.RequestType -}}
    {{- $rpcRspType := $method.ResponseType -}}
	{{ with .LeadingComments }}// {{$rpcName}} {{.}}{{ end }}
	{{$rpcName}}(ctx context.Context, req *{{$rpcReqType}},rsp *{{$rpcRspType}}) (err error) {{ with .TrailingComments}}// {{.}}{{ end }}
{{ end -}}
}

{{range $index, $method := $service.RPC -}}
{{- $rpcName := $method.Name | camelcase -}}
{{- $rpcReqType := $method.RequestType -}}
{{- $rpcRspType := $method.ResponseType -}}
func {{$svrNameCamelCase}}Service_{{$rpcName}}_Handler(svr interface{}, ctx context.Context, f server.FilterFunc) (rspbody interface{}, err error) {

    req := &{{$rpcReqType}}{}
	rsp := &{{$rpcRspType}}{}
	filters, err := f(req)
    if err != nil {
    	return nil, err
    }
	handleFunc := func(ctx context.Context, reqbody interface{}, rspbody interface{}) error {
	    return svr.({{$svrNameCamelCase}}Service).{{$rpcName}}(ctx, reqbody.(*{{$rpcReqType}}), rspbody.(*{{$rpcRspType}}))
	}

	err = filters.Handle(ctx, req, rsp, handleFunc)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

{{end -}}

// {{$svrNameCamelCase}}Server_ServiceDesc descriptor for server.RegisterService
var {{$svrNameCamelCase}}Server_ServiceDesc = server.ServiceDesc {
    ServiceName: "{{$pkgName}}.{{$svrName}}",
    HandlerType: ((*{{$svrNameCamelCase}}Service)(nil)),
    Methods: []server.Method{
        {{- range $service.RPC}}
        {{- $rpcName := .Name | camelcase -}}
	        {Name: "{{.FullyQualifiedCmd}}", Func: {{$svrNameCamelCase}}Service_{{$rpcName}}_Handler},
        {{- end}}
    },
}

// Register{{$svrNameCamelCase}}Service register service
func Register{{$svrNameCamelCase}}Service(s server.Service, svr {{$svrNameCamelCase}}Service) {
	s.Register(&{{$svrNameCamelCase}}Server_ServiceDesc, svr)
}

{{ end }}

/* ************************************ Client Definition ************************************ */

{{ range $service := .Services }}
{{ $svrNameCamelCase := $service.Name | camelcase }}
// {{$svrNameCamelCase}}ClientProxy defines service client proxy
type {{$svrNameCamelCase}}ClientProxy interface {
	{{ range $rpc := $service.RPC}}
   	{{- $rpcName := .Name | camelcase -}}
    {{- $rpcReqType := $rpc.RequestType -}}
    {{- $rpcRspType := $rpc.ResponseType -}}
   	{{ with .LeadingComments }}// {{$rpcName}} {{.}}{{ end }}
	{{$rpcName}}(ctx context.Context, req *{{$rpcReqType}}, opts ...client.Option) (rsp *{{$rpcRspType}}, err error) {{ with .TrailingComments }}// {{.}}{{ end }}
{{ end -}}
}

type {{$svrNameCamelCase|untitle}}ClientProxyImpl struct{
	client client.Client
	opts []client.Option
}

func New{{$svrNameCamelCase}}ClientProxy(opts...client.Option) {{$svrNameCamelCase}}ClientProxy {
	return &{{$svrNameCamelCase|untitle}}ClientProxyImpl {client: client.DefaultClient, opts: opts}
}

{{range $index, $method := $service.RPC}}
{{- $rpcName := $method.Name | camelcase -}}
{{- $rpcReqType := $method.RequestType -}}
{{- $rpcRspType := $method.ResponseType -}}
func (c *{{$svrNameCamelCase|untitle}}ClientProxyImpl) {{$rpcName}}(ctx context.Context, req *{{$rpcReqType}}, opts ...client.Option) (rsp *{{$rpcRspType}}, err error) {

	ctx, msg := codec.WithCloneMessage(ctx)

	msg.WithClientRPCName({{$svrNameCamelCase}}Server_ServiceDesc.Methods[{{$index}}].Name)
	msg.WithCalleeServiceName({{$svrNameCamelCase}}Server_ServiceDesc.ServiceName)

	callopts := make([]client.Option, 0, len(c.opts)+len(opts))
	callopts = append(callopts, c.opts...)
	callopts = append(callopts, opts...)
	rsp = &{{$rpcRspType}}{}

	err = c.client.Invoke(ctx, req, rsp, callopts...)
	if err != nil {
	    return nil, err
	}

	return rsp, nil
}
{{end}}

{{ end }}
