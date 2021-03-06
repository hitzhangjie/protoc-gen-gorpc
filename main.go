// Go support for Protocol Buffers - Google's data interchange format
//
// Copyright 2010 The Go Authors.  All rights reserved.
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

// protoc-gen-gorpc is a plugin for the Google protocol buffer compiler to generate
// Go code.  Run it by building this program and putting it in your path with
// the name
// 	protoc-gen-gorpc
// That word 'go' at the end becomes part of the option string set for the
// protocol compiler, so once the protocol compiler (protoc) is installed
// you can run
// 	protoc --gorpc_out=output_directory input_directory/file.proto
// to generate Go bindings for the protocol defined by file.proto.
// With that input, the output will be written to
// 	output_directory/file.pb.go
//
// The generated code is documented in the package comment for
// the library.
//
// See the README and documentation for protocol buffers to learn more:
// 	https://developers.google.com/protocol-buffers/
package main

import (
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/hitzhangjie/protoc-gen-gorpc/generator"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)

	fout, err := os.OpenFile("/tmp/protoc-gen-gorpc.log", os.O_CREATE|os.O_RDWR, os.FileMode(0666))
	if err != nil {
		panic(err)
	}
	log.SetOutput(fout)
}

func main() {

	//============================================================
	// 0. prepare the generator

	// Begin by allocating a generator. The request and response structures are stored there
	// so we can do error handling easily - the response structure contains the field to
	// report failure.
	g := generator.New()

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		g.Error(err, "reading input")
	}

	if err := proto.Unmarshal(data, g.Request); err != nil {
		g.Error(err, "parsing input proto")
	}

	if len(g.Request.FileToGenerate) == 0 {
		g.Fail("no files to generate")
	}

	g.CommandLineParameters(g.Request.GetParameter())

	// if want to debug this protoc plugin, we'd better let the debugger attach to it before
	// the execution of plugin logic.
	if v, ok := g.Param["debug"]; ok && len(v) != 0 {
		if debug, err := strconv.ParseBool(v); err == nil {
			if debug {
				log.Printf("protoc-gen-gorpc pid: %d, ready to pause", os.Getpid())
			}
			for debug {
				runtime.Gosched()
			}
		}
	}

	// Create a wrapped version of the Descriptors and EnumDescriptors that
	// point to the file that defines them.
	g.WrapTypes()

	g.SetPackageNames()
	g.BuildTypeNameMap()

	//============================================================
	// 1. generate *.pb.go，here the grpc logic is removed

	g.GenerateAllFiles()

	//============================================================
	// 2. process go template files registered in gotpl.GoRPCTemplates

	err = g.GenerateTplFiles()
	if err != nil {
		log.Printf("failed to process template files, err: %v", err)
	}

	//Deprecated: run `go fmt`
	//
	//dir, err := os.Getwd()
	//if err != nil {
	//	log.Printf("failed to get current directory: %v", err)
	//}
	//
	//cmd := exec.Command("go", "fmt", dir)
	//info, err := cmd.CombinedOutput()
	//if err != nil {
	//	log.Printf("failed to gofmt directory, err: %v, info: %s", err, string(info))
	//}
	//log.Printf("run gofmt ok:\n%s", string(info))

	//============================================================
	// 3. Send back the results.

	data, err = proto.Marshal(g.Response)
	if err != nil {
		g.Error(err, "failed to marshal output proto")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		g.Error(err, "failed to write output proto")
	}

}
