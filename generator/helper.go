package generator

import (
	"github.com/hitzhangjie/protoc-gen-gorpc/descriptor"
	"github.com/hitzhangjie/protoc-gen-gorpc/gorpc"
)

// BuildFileDescriptor 构建一个更简单的FileDescriptor对象，指导代码生成
func BuildFileDescriptor(fd *FileDescriptor) (nfd *gorpc.FileDescriptor, err error) {

	opts, err := buildOptions(fd.GetOptions())
	if err != nil {
		return nil, err
	}

	nfd = &gorpc.FileDescriptor{
		PackageName:        fd.GetPackage(),
		Imports:            nil, // fixme
		FileOptions:        opts,
		Services:           []*gorpc.ServiceDescriptor{},
		Dependencies:       nil, // fixme
		ImportPathMappings: nil, // fixme
		PkgPkgMappings:     nil, // fixme
	}

	for _, s := range fd.Service {
		srv := &gorpc.ServiceDescriptor{
			Name: s.GetName(),
			RPC:  []*gorpc.RPCDescriptor{},
		}
		for _, m := range s.Method {
			rpc := &gorpc.RPCDescriptor{
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
		nfd.Services = append(nfd.Services, srv)
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
