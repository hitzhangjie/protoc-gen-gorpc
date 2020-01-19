// Go support for Protocol Buffers - Google's data interchange format
//
// Copyright 2015 The Go Authors.  All rights reserved.
// https://github.com/golang/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Package gorpc outputs gorpc service descriptions in Go code, besides,
// it also outputs gorpc server project and client stubs.
//
// It runs as a plugin for the Go protocol buffer compiler plugin.
// It is linked in to protoc-gen-gorpc.
package gorpc

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/hitzhangjie/protoc-gen-gorpc/utils/fs"

	pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
)

// generatedCodeVersion indicates a version of the generated code.
// It is incremented whenever an incompatibility between the generated code and
// the gorpc package is introduced; the generated code references
// a constant, gorpc.SupportPackageIsVersionN (where N is generatedCodeVersion).
const generatedCodeVersion = 4

// Paths for packages used by code generated in this file,
// relative to the import_prefix of the generator.Generator.
const (
	contextPkgPath = "context"

	// fixme, remove this
	gorpcPkgPath  = "google.golang.org/gorpc"
	codePkgPath   = "google.golang.org/gorpc/codes"
	statusPkgPath = "google.golang.org/gorpc/status"
)

func init() {
	generator.RegisterPlugin(new(gorpc))
}

// gorpc is an implementation of the Go protocol buffer compiler's
// plugin architecture.  It generates bindings for gorpc support.
type gorpc struct {
	gen *generator.Generator
}

// Name returns the name of this plugin, "gorpc".
func (g *gorpc) Name() string {
	return "gorpc"
}

// The names for packages imported in the generated code.
// They may vary from the final path component of the import path
// if the name is used by other packages.
var (
	contextPkg string
	gorpcPkg   string
)

// Init initializes the plugin.
func (g *gorpc) Init(gen *generator.Generator) {
	g.gen = gen
}

// Given a type name defined in a .proto, return its object.
// Also record that we're using it, to guarantee the associated import.
func (g *gorpc) objectNamed(name string) generator.Object {
	g.gen.RecordTypeUse(name)
	return g.gen.ObjectNamed(name)
}

// Given a type name defined in a .proto, return its name as we will print it.
func (g *gorpc) typeName(str string) string {
	return g.gen.TypeName(g.objectNamed(str))
}

// P forwards to g.gen.P.
func (g *gorpc) P(args ...interface{}) { g.gen.P(args...) }

// Generate generates code for the services in the given file.
func (g *gorpc) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}

	// todo build the FileDescriptor
	nfd, err := buildFileDescriptor(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "buildFileDescriptor error: %v", err)
		os.Exit(1)
	}

	// todo run go template to generate template
	root := "/Users/zhangjie/Github/protoc-gen-gorpc/install"
	tmpdir := os.TempDir()

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		// 检查要不要处理当前文件
		if err != nil {
			return err
		}

		if path == "." || path == ".." {
			return nil
		}

		var target string

		// 新生成文件目录结构，与模板路径保持一样的结构
		if rel, err := filepath.Rel(path, root); err != nil {
			return err
		} else {
			target = filepath.Join(tmpdir, rel)
		}

		// 如果是文件，且为go模板文件，执行go模板引擎生成新文件
		if !info.IsDir() {

			// 非模板文件，直接copy
			if !strings.HasSuffix(path, ".tpl") {
				fs.Copy(path, target)
			}

			// 模板文件，执行模板处理引擎
			if err = processTemplateFile(path, target, nfd); err != nil {
				return err
			}
			return nil
		}

		if err = os.MkdirAll(target, os.ModePerm); err != nil {
			return err
		}
		return nil
	})

}

func processTemplateFile(inFile, outFile string, nfd *FileDescriptor) error {

	baseName := filepath.Base(inFile)

	var (
		instance *template.Template
		err      error
		fout     *os.File
	)

	if funcMap == nil {
		instance, err = template.New(baseName).ParseFiles(inFile)
	} else {
		instance, err = template.New(baseName).Funcs(funcMap).ParseFiles(inFile)
	}

	if err != nil {
		return err
	}

	if fout, err = os.Create(outFile); err != nil {
		return err
	}
	defer fout.Close()

	if err = instance.Execute(fout, nfd); err != nil {
		return err
	}
	return nil
}

// GenerateImports generates the import declaration for this file.
func (g *gorpc) GenerateImports(file *generator.FileDescriptor) {
}

// reservedClientName records whether a client name is reserved on the client side.
var reservedClientName = map[string]bool{
	// TODO: do we need any in gRPC?
}

func unexport(s string) string { return strings.ToLower(s[:1]) + s[1:] }

// deprecationComment is the standard comment added to deprecated
// messages, fields, enums, and enum values.
var deprecationComment = "// Deprecated: Do not use."

// generateService generates all the code for the named service.
func (g *gorpc) generateService(file *generator.FileDescriptor, service *pb.ServiceDescriptorProto, index int) {
	path := fmt.Sprintf("6,%d", index) // 6 means service.

	origServName := service.GetName()
	fullServName := origServName
	if pkg := file.GetPackage(); pkg != "" {
		fullServName = pkg + "." + fullServName
	}
	servName := generator.CamelCase(origServName)
	deprecated := service.GetOptions().GetDeprecated()

	g.P()
	g.P(fmt.Sprintf(`// %sClient is the client API for %s service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/gorpc#ClientConn.NewStream.`, servName, servName))

	// Client interface.
	if deprecated {
		g.P("//")
		g.P(deprecationComment)
	}
	g.P("type ", servName, "Client interface {")
	for i, method := range service.Method {
		g.gen.PrintComments(fmt.Sprintf("%s,2,%d", path, i)) // 2 means method in a service.
		if method.GetOptions().GetDeprecated() {
			g.P("//")
			g.P(deprecationComment)
		}
		g.P(g.generateClientSignature(servName, method))
	}
	g.P("}")
	g.P()

	// Client structure.
	g.P("type ", unexport(servName), "Client struct {")
	g.P("cc *", gorpcPkg, ".ClientConn")
	g.P("}")
	g.P()

	// NewClient factory.
	if deprecated {
		g.P(deprecationComment)
	}
	g.P("func New", servName, "Client (cc *", gorpcPkg, ".ClientConn) ", servName, "Client {")
	g.P("return &", unexport(servName), "Client{cc}")
	g.P("}")
	g.P()

	var methodIndex, streamIndex int
	serviceDescVar := "_" + servName + "_serviceDesc"
	// Client method implementations.
	for _, method := range service.Method {
		var descExpr string
		if !method.GetServerStreaming() && !method.GetClientStreaming() {
			// Unary RPC method
			descExpr = fmt.Sprintf("&%s.Methods[%d]", serviceDescVar, methodIndex)
			methodIndex++
		} else {
			// Streaming RPC method
			descExpr = fmt.Sprintf("&%s.Streams[%d]", serviceDescVar, streamIndex)
			streamIndex++
		}
		g.generateClientMethod(servName, fullServName, serviceDescVar, method, descExpr)
	}

	// Server interface.
	serverType := servName + "Server"
	g.P("// ", serverType, " is the server API for ", servName, " service.")
	if deprecated {
		g.P("//")
		g.P(deprecationComment)
	}
	g.P("type ", serverType, " interface {")
	for i, method := range service.Method {
		g.gen.PrintComments(fmt.Sprintf("%s,2,%d", path, i)) // 2 means method in a service.
		if method.GetOptions().GetDeprecated() {
			g.P("//")
			g.P(deprecationComment)
		}
		g.P(g.generateServerSignature(servName, method))
	}
	g.P("}")
	g.P()

	// Server Unimplemented struct for forward compatibility.
	if deprecated {
		g.P(deprecationComment)
	}
	g.generateUnimplementedServer(servName, service)

	// Server registration.
	if deprecated {
		g.P(deprecationComment)
	}
	g.P("func Register", servName, "Server(s *", gorpcPkg, ".Server, srv ", serverType, ") {")
	g.P("s.RegisterService(&", serviceDescVar, `, srv)`)
	g.P("}")
	g.P()

	// Server handler implementations.
	var handlerNames []string
	for _, method := range service.Method {
		hname := g.generateServerMethod(servName, fullServName, method)
		handlerNames = append(handlerNames, hname)
	}

	// Service descriptor.
	g.P("var ", serviceDescVar, " = ", gorpcPkg, ".ServiceDesc {")
	g.P("ServiceName: ", strconv.Quote(fullServName), ",")
	g.P("HandlerType: (*", serverType, ")(nil),")
	g.P("Methods: []", gorpcPkg, ".MethodDesc{")
	for i, method := range service.Method {
		if method.GetServerStreaming() || method.GetClientStreaming() {
			continue
		}
		g.P("{")
		g.P("MethodName: ", strconv.Quote(method.GetName()), ",")
		g.P("Handler: ", handlerNames[i], ",")
		g.P("},")
	}
	g.P("},")
	g.P("Streams: []", gorpcPkg, ".StreamDesc{")
	for i, method := range service.Method {
		if !method.GetServerStreaming() && !method.GetClientStreaming() {
			continue
		}
		g.P("{")
		g.P("StreamName: ", strconv.Quote(method.GetName()), ",")
		g.P("Handler: ", handlerNames[i], ",")
		if method.GetServerStreaming() {
			g.P("ServerStreams: true,")
		}
		if method.GetClientStreaming() {
			g.P("ClientStreams: true,")
		}
		g.P("},")
	}
	g.P("},")
	g.P("Metadata: \"", file.GetName(), "\",")
	g.P("}")
	g.P()
}

// generateUnimplementedServer creates the unimplemented server struct
func (g *gorpc) generateUnimplementedServer(servName string, service *pb.ServiceDescriptorProto) {
	serverType := servName + "Server"
	g.P("// Unimplemented", serverType, " can be embedded to have forward compatible implementations.")
	g.P("type Unimplemented", serverType, " struct {")
	g.P("}")
	g.P()
	// Unimplemented<service_name>Server's concrete methods
	for _, method := range service.Method {
		g.generateServerMethodConcrete(servName, method)
	}
	g.P()
}

// generateServerMethodConcrete returns unimplemented methods which ensure forward compatibility
func (g *gorpc) generateServerMethodConcrete(servName string, method *pb.MethodDescriptorProto) {
	header := g.generateServerSignatureWithParamNames(servName, method)
	g.P("func (*Unimplemented", servName, "Server) ", header, " {")
	var nilArg string
	if !method.GetServerStreaming() && !method.GetClientStreaming() {
		nilArg = "nil, "
	}
	methName := generator.CamelCase(method.GetName())
	statusPkg := string(g.gen.AddImport(statusPkgPath))
	codePkg := string(g.gen.AddImport(codePkgPath))
	g.P("return ", nilArg, statusPkg, `.Errorf(`, codePkg, `.Unimplemented, "method `, methName, ` not implemented")`)
	g.P("}")
}

// generateClientSignature returns the client-side signature for a method.
func (g *gorpc) generateClientSignature(servName string, method *pb.MethodDescriptorProto) string {
	origMethName := method.GetName()
	methName := generator.CamelCase(origMethName)
	if reservedClientName[methName] {
		methName += "_"
	}
	reqArg := ", in *" + g.typeName(method.GetInputType())
	if method.GetClientStreaming() {
		reqArg = ""
	}
	respName := "*" + g.typeName(method.GetOutputType())
	if method.GetServerStreaming() || method.GetClientStreaming() {
		respName = servName + "_" + generator.CamelCase(origMethName) + "Client"
	}
	return fmt.Sprintf("%s(ctx %s.Context%s, opts ...%s.CallOption) (%s, error)", methName, contextPkg, reqArg, gorpcPkg, respName)
}

func (g *gorpc) generateClientMethod(servName, fullServName, serviceDescVar string, method *pb.MethodDescriptorProto, descExpr string) {
	sname := fmt.Sprintf("/%s/%s", fullServName, method.GetName())
	methName := generator.CamelCase(method.GetName())
	inType := g.typeName(method.GetInputType())
	outType := g.typeName(method.GetOutputType())

	if method.GetOptions().GetDeprecated() {
		g.P(deprecationComment)
	}
	g.P("func (c *", unexport(servName), "Client) ", g.generateClientSignature(servName, method), "{")
	if !method.GetServerStreaming() && !method.GetClientStreaming() {
		g.P("out := new(", outType, ")")
		// TODO: Pass descExpr to Invoke.
		g.P(`err := c.cc.Invoke(ctx, "`, sname, `", in, out, opts...)`)
		g.P("if err != nil { return nil, err }")
		g.P("return out, nil")
		g.P("}")
		g.P()
		return
	}
	streamType := unexport(servName) + methName + "Client"
	g.P("stream, err := c.cc.NewStream(ctx, ", descExpr, `, "`, sname, `", opts...)`)
	g.P("if err != nil { return nil, err }")
	g.P("x := &", streamType, "{stream}")
	if !method.GetClientStreaming() {
		g.P("if err := x.ClientStream.SendMsg(in); err != nil { return nil, err }")
		g.P("if err := x.ClientStream.CloseSend(); err != nil { return nil, err }")
	}
	g.P("return x, nil")
	g.P("}")
	g.P()

	genSend := method.GetClientStreaming()
	genRecv := method.GetServerStreaming()
	genCloseAndRecv := !method.GetServerStreaming()

	// Stream auxiliary types and methods.
	g.P("type ", servName, "_", methName, "Client interface {")
	if genSend {
		g.P("Send(*", inType, ") error")
	}
	if genRecv {
		g.P("Recv() (*", outType, ", error)")
	}
	if genCloseAndRecv {
		g.P("CloseAndRecv() (*", outType, ", error)")
	}
	g.P(gorpcPkg, ".ClientStream")
	g.P("}")
	g.P()

	g.P("type ", streamType, " struct {")
	g.P(gorpcPkg, ".ClientStream")
	g.P("}")
	g.P()

	if genSend {
		g.P("func (x *", streamType, ") Send(m *", inType, ") error {")
		g.P("return x.ClientStream.SendMsg(m)")
		g.P("}")
		g.P()
	}
	if genRecv {
		g.P("func (x *", streamType, ") Recv() (*", outType, ", error) {")
		g.P("m := new(", outType, ")")
		g.P("if err := x.ClientStream.RecvMsg(m); err != nil { return nil, err }")
		g.P("return m, nil")
		g.P("}")
		g.P()
	}
	if genCloseAndRecv {
		g.P("func (x *", streamType, ") CloseAndRecv() (*", outType, ", error) {")
		g.P("if err := x.ClientStream.CloseSend(); err != nil { return nil, err }")
		g.P("m := new(", outType, ")")
		g.P("if err := x.ClientStream.RecvMsg(m); err != nil { return nil, err }")
		g.P("return m, nil")
		g.P("}")
		g.P()
	}
}

// generateServerSignatureWithParamNames returns the server-side signature for a method with parameter names.
func (g *gorpc) generateServerSignatureWithParamNames(servName string, method *pb.MethodDescriptorProto) string {
	origMethName := method.GetName()
	methName := generator.CamelCase(origMethName)
	if reservedClientName[methName] {
		methName += "_"
	}

	var reqArgs []string
	ret := "error"
	if !method.GetServerStreaming() && !method.GetClientStreaming() {
		reqArgs = append(reqArgs, "ctx "+contextPkg+".Context")
		ret = "(*" + g.typeName(method.GetOutputType()) + ", error)"
	}
	if !method.GetClientStreaming() {
		reqArgs = append(reqArgs, "req *"+g.typeName(method.GetInputType()))
	}
	if method.GetServerStreaming() || method.GetClientStreaming() {
		reqArgs = append(reqArgs, "srv "+servName+"_"+generator.CamelCase(origMethName)+"Server")
	}

	return methName + "(" + strings.Join(reqArgs, ", ") + ") " + ret
}

// generateServerSignature returns the server-side signature for a method.
func (g *gorpc) generateServerSignature(servName string, method *pb.MethodDescriptorProto) string {
	origMethName := method.GetName()
	methName := generator.CamelCase(origMethName)
	if reservedClientName[methName] {
		methName += "_"
	}

	var reqArgs []string
	ret := "error"
	if !method.GetServerStreaming() && !method.GetClientStreaming() {
		reqArgs = append(reqArgs, contextPkg+".Context")
		ret = "(*" + g.typeName(method.GetOutputType()) + ", error)"
	}
	if !method.GetClientStreaming() {
		reqArgs = append(reqArgs, "*"+g.typeName(method.GetInputType()))
	}
	if method.GetServerStreaming() || method.GetClientStreaming() {
		reqArgs = append(reqArgs, servName+"_"+generator.CamelCase(origMethName)+"Server")
	}

	return methName + "(" + strings.Join(reqArgs, ", ") + ") " + ret
}

func (g *gorpc) generateServerMethod(servName, fullServName string, method *pb.MethodDescriptorProto) string {
	methName := generator.CamelCase(method.GetName())
	hname := fmt.Sprintf("_%s_%s_Handler", servName, methName)
	inType := g.typeName(method.GetInputType())
	outType := g.typeName(method.GetOutputType())

	if !method.GetServerStreaming() && !method.GetClientStreaming() {
		g.P("func ", hname, "(srv interface{}, ctx ", contextPkg, ".Context, dec func(interface{}) error, interceptor ", gorpcPkg, ".UnaryServerInterceptor) (interface{}, error) {")
		g.P("in := new(", inType, ")")
		g.P("if err := dec(in); err != nil { return nil, err }")
		g.P("if interceptor == nil { return srv.(", servName, "Server).", methName, "(ctx, in) }")
		g.P("info := &", gorpcPkg, ".UnaryServerInfo{")
		g.P("Server: srv,")
		g.P("FullMethod: ", strconv.Quote(fmt.Sprintf("/%s/%s", fullServName, methName)), ",")
		g.P("}")
		g.P("handler := func(ctx ", contextPkg, ".Context, req interface{}) (interface{}, error) {")
		g.P("return srv.(", servName, "Server).", methName, "(ctx, req.(*", inType, "))")
		g.P("}")
		g.P("return interceptor(ctx, in, info, handler)")
		g.P("}")
		g.P()
		return hname
	}
	streamType := unexport(servName) + methName + "Server"
	g.P("func ", hname, "(srv interface{}, stream ", gorpcPkg, ".ServerStream) error {")
	if !method.GetClientStreaming() {
		g.P("m := new(", inType, ")")
		g.P("if err := stream.RecvMsg(m); err != nil { return err }")
		g.P("return srv.(", servName, "Server).", methName, "(m, &", streamType, "{stream})")
	} else {
		g.P("return srv.(", servName, "Server).", methName, "(&", streamType, "{stream})")
	}
	g.P("}")
	g.P()

	genSend := method.GetServerStreaming()
	genSendAndClose := !method.GetServerStreaming()
	genRecv := method.GetClientStreaming()

	// Stream auxiliary types and methods.
	g.P("type ", servName, "_", methName, "Server interface {")
	if genSend {
		g.P("Send(*", outType, ") error")
	}
	if genSendAndClose {
		g.P("SendAndClose(*", outType, ") error")
	}
	if genRecv {
		g.P("Recv() (*", inType, ", error)")
	}
	g.P(gorpcPkg, ".ServerStream")
	g.P("}")
	g.P()

	g.P("type ", streamType, " struct {")
	g.P(gorpcPkg, ".ServerStream")
	g.P("}")
	g.P()

	if genSend {
		g.P("func (x *", streamType, ") Send(m *", outType, ") error {")
		g.P("return x.ServerStream.SendMsg(m)")
		g.P("}")
		g.P()
	}
	if genSendAndClose {
		g.P("func (x *", streamType, ") SendAndClose(m *", outType, ") error {")
		g.P("return x.ServerStream.SendMsg(m)")
		g.P("}")
		g.P()
	}
	if genRecv {
		g.P("func (x *", streamType, ") Recv() (*", inType, ", error) {")
		g.P("m := new(", inType, ")")
		g.P("if err := x.ServerStream.RecvMsg(m); err != nil { return nil, err }")
		g.P("return m, nil")
		g.P("}")
		g.P()
	}

	return hname
}
