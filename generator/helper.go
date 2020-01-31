package generator

import (
	"github.com/hitzhangjie/protoc-gen-gorpc/gorpc"
)

// convertFileDescriptor convert google/protobuf/.../FileDescriptor to gorpc.FileDescriptor,
// which is much simpler to reference in go template file.
func convertFileDescriptor(fd *FileDescriptor) (nfd *gorpc.FileDescriptor, err error) {

	// package name
	packageName := fd.GetPackage()
	if opts := fd.GetOptions(); opts != nil {
		if v := opts.GetGoPackage(); len(v) != 0 {
			packageName = v
		}
	}

	nfd = &gorpc.FileDescriptor{
		PackageName:        packageName,
		Imports:            nil, // fixme
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
