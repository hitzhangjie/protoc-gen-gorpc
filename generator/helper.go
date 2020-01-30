package gorpc

import (
	"github.com/hitzhangjie/protoc-gen-gorpc/descriptor"
	"github.com/hitzhangjie/protoc-gen-gorpc/generator"
)

// BuildFileDescriptor 构建一个更简单的FileDescriptor对象，指导代码生成
func BuildFileDescriptor(fd *generator.FileDescriptor) (nfd *FileDescriptor, err error) {

	opts, err := buildOptions(fd.GetOptions())
	if err != nil {
		return nil, err
	}

	nfd = &FileDescriptor{
		PackageName:        fd.GetPackage(),
		Imports:            nil, // fixme
		FileOptions:        opts,
		Services:           []*ServiceDescriptor{},
		Dependencies:       nil, // fixme
		ImportPathMappings: nil, // fixme
		pkgPkgMappings:     nil, // fixme
	}

	for _, s := range fd.Service {
		srv := &ServiceDescriptor{
			Name: s.GetName(),
			RPC:  []*RPCDescriptor{},
		}
		for _, m := range s.Method {
			rpc := &RPCDescriptor{
				Name:              m.GetName(),
				Cmd:               m.GetName(),
				FullyQualifiedCmd: m.GetName(),
				RequestType:       m.GetInputType(),
				ResponseType:      m.GetOutputType(),
				LeadingComments:   "", // fixme
				TrailingComments:  "", // fixme
			}
			srv.RPC = append(srv.RPC, rpc)
		}
	}

	return nfd, nil
}

// buildOptions 将常见的影响package名的FileOptions转换为map
func buildOptions(fopts *descriptor.FileOptions) (map[string]interface{}, error) {

	if fopts == nil {
		return nil, nil
	}

	m := map[string]interface{}{}

	if v := fopts.GetGoPackage(); len(v) != 0 {
		m["go_package"] = v
	}
	//if v := fopts.GetJavaPackage(); len(v) != 0 {
	//	m["java_package"] = v
	//}

	return m, nil
}
