package gorpc

// FileDescriptor 文件作用域相关的描述信息
type FileDescriptor struct {
	ProtoFile string
	FilePath  string

	// PackageName package name in generated *.go, which is determined by
	// - `option go_package=`,
	// - or `package directive` if `option go_package=` not specified.
	PackageName string // pb文件package diretive确定的包名

	Services []*ServiceDescriptor // 支持多service，目前只处理第一个service

	Dependencies       map[string]string // 依赖(imported)的pb文件对应的正确包名（考虑了fileoptions如go_package等的影响）
	ImportPathMappings map[string]string // pb文件package(package diretive)到正确导入路径的关系(考虑了fileoptions如go_package的响应)
	PkgPkgMappings     map[string]string // pb文件package(package directive)到正确包名的映射关系(考虑了fileoptions如go_package的影响）

	// fixme remove
	Pb2ValidGoPkg  map[string]string // k=pb文件名, v=protoc处理后package名
	Pkg2ValidGoPkg map[string]string // k=pb文件package directive, v=protoc处理后package名
	Pkg2ImportPath map[string]string // k=pb文件package directive, v=go代码中对应importpath
	Pb2ImportPath  map[string]string // k=pb文件名，v=go代码中对应importpath

	// Deprecated
	Imports []string // pb文件可能import其他pb文件，登记rpc请求、响应中引用的package (package diretive确定的包名)
	// Deprecated
	FileOptions map[string]interface{} // fileoptions，如go_package, java_package等
}

// ServiceDescriptor service作用域相关的描述信息
type ServiceDescriptor struct {
	Name string           // 服务名称
	RPC  []*RPCDescriptor // rpc接口定义
}

// RPCDescriptor rpc作用域相关的描述信息
//
// RequestType由于涉及到
type RPCDescriptor struct {
	Name              string // RPC方法名
	Cmd               string // RPC命令字
	FullyQualifiedCmd string // 完整的RPC命令字，用于ServiceDesc、client请求时命令字
	RequestType       string // RPC请求消息类型，包含package，比如package_a.TypeA
	ResponseType      string // RPC响应消息类型，包含package，比如package_b.TypeB
	LeadingComments   string // RPC前置注释信息
	TrailingComments  string // RPC后置注释信息
}
