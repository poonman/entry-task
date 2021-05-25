package dora

import (
	"path"
	"strconv"
	"strings"

	pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/poonman/entry-task/dora/misc/util/protoc-gen-dora/generator"
)

// Paths for packages used by code generated in this file,
// relative to the import_prefix of the generator.Generator.
const (
	statusPkgPath  = "github.com/poonman/entry-task/dora/status"
	contextPkgPath = "context"
	clientPkgPath  = "github.com/poonman/entry-task/dora/client"
	serverPkgPath  = "github.com/poonman/entry-task/dora/server"
)

func init() {
	generator.RegisterPlugin(new(doraPlugin))
}

// doraPlugin is an implementation of the Go protocol buffer compiler's
// plugin architecture.  It generates bindings for go-doraPlugin support.
type doraPlugin struct {
	gen *generator.Generator
}

// Name returns the name of this plugin, "doraPlugin".
func (g *doraPlugin) Name() string {
	return "doraPlugin"
}

// The names for packages imported in the generated code.
// They may vary from the final path component of the import path
// if the name is used by other packages.
var (
	statusPkg  string
	contextPkg string
	clientPkg  string
	serverPkg  string
	pkgImports map[generator.GoPackageName]bool
)

// Init initializes the plugin.
func (g *doraPlugin) Init(gen *generator.Generator) {
	g.gen = gen
	statusPkg = generator.RegisterUniquePackageName("status", nil)
	contextPkg = generator.RegisterUniquePackageName("context", nil)
	clientPkg = generator.RegisterUniquePackageName("client", nil)
	serverPkg = generator.RegisterUniquePackageName("server", nil)
}

// Given a type name defined in a .proto, return its object.
// Also record that we're using it, to guarantee the associated import.
func (g *doraPlugin) objectNamed(name string) generator.Object {
	g.gen.RecordTypeUse(name)
	return g.gen.ObjectNamed(name)
}

// Given a type name defined in a .proto, return its name as we will print it.
func (g *doraPlugin) typeName(str string) string {
	return g.gen.TypeName(g.objectNamed(str))
}

// P forwards to g.gen.P.
func (g *doraPlugin) P(args ...interface{}) { g.gen.P(args...) }

// Generate generates code for the services in the given file.
func (g *doraPlugin) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
	//g.P("// Reference imports to suppress errors if they are not otherwise used.")
	//g.P("var _ ", statusPkg, ".Endpoint")
	//g.P("var _ ", contextPkg, ".Context")
	//g.P("var _ ", clientPkg, ".Option")
	//g.P("var _ ", serverPkg, ".Option")
	//g.P()

	for i, service := range file.FileDescriptorProto.Service {
		g.generateDoraService(file, service, i)
	}
}

// GenerateImports generates the import declaration for this file.
func (g *doraPlugin) GenerateImports(file *generator.FileDescriptor, imports map[generator.GoImportPath]generator.GoPackageName) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
	g.P("import (")
	g.P(statusPkg, " ", strconv.Quote(path.Join(g.gen.ImportPrefix, statusPkgPath)))
	g.P(contextPkg, " ", strconv.Quote(path.Join(g.gen.ImportPrefix, contextPkgPath)))
	g.P(clientPkg, " ", strconv.Quote(path.Join(g.gen.ImportPrefix, clientPkgPath)))
	g.P(serverPkg, " ", strconv.Quote(path.Join(g.gen.ImportPrefix, serverPkgPath)))
	g.P(")")
	g.P()

	// We need to keep track of imported packages to make sure we don't produce
	// a name collision when generating types.
	pkgImports = make(map[generator.GoPackageName]bool)
	for _, name := range imports {
		pkgImports[name] = true
	}
}

// reservedClientName records whether a client name is reserved on the client side.
var reservedClientName = map[string]bool{
	// TODO: do we need any in go-doraPlugin?
}

func unexport(s string) string {
	if len(s) == 0 {
		return ""
	}
	name := strings.ToLower(s[:1]) + s[1:]
	if pkgImports[generator.GoPackageName(name)] {
		return name + "_"
	}
	return name
}

func (g *doraPlugin) generateDoraService(file *generator.FileDescriptor, service *pb.ServiceDescriptorProto, index int) {
	//path := fmt.Sprintf("6,%d", index) // 6 means service.

	origServName := service.GetName()
	//serviceName := strings.ToLower(service.GetName())
	//if pkg := file.GetPackage(); pkg != "" {
	//	serviceName = pkg
	//}
	servName := generator.CamelCase(origServName)
	servAlias := servName + "Service"

	// strip suffix
	if strings.HasSuffix(servAlias, "ServiceService") {
		servAlias = strings.TrimSuffix(servAlias, "Service")
	}

	g.P("type ", servName, "Client interface {")
	for _, method := range service.Method {
		g.P(method.GetName(), "(ctx context.Context, in *", method.GetInputType(), ") (out *", method.GetOutputType(), ", err error)")
	}
	g.P("}")

	g.P()

	lowerServName := unexport(servName)
	g.P("type ", lowerServName, "Client struct {")
	g.P("cc client.Invoker")
	g.P("}")

	g.P()

	for _, method := range service.Method {
		g.P("func (c *", lowerServName, "Client) ", method.GetName(),
			"(ctx context.Context, in *", method.GetInputType(), ") (out *", method.GetOutputType(), ", err error) {")
		g.P("out = &", method.GetOutputType(), "{}")
		g.P()
		g.P(`err = c.cc.Invoke(ctx, "`, method.GetName(), `", in, out)`)
		g.P("if err != nil {")
		g.P("return nil, err")
		g.P("}")
		g.P("return out, nil")
	}

	// server
	g.P("type ", servName, "Server interface {")
	for _, method := range service.Method {
		g.P(method.GetName(), "(ctx context.Context, in *", method.GetInputType(), ") (out *", method.GetOutputType(), ", err error)")
	}
	g.P("mustEmbedUnimplemented", servName, "Server()")
	g.P("}")

	g.P()

	g.P("type Unimplemented", servName, "Server struct {")
	g.P("}")

	g.P()

	for _, method := range service.Method {
		g.P("func (Unimplemented", servName, "Server) ", method.GetName(), "(context.Context, *", method.GetInputType(), ") (*", method.GetOutputType(), ", error) {")
		g.P(`return nil, status.New(status.Unimplemented, "server is unimplemented")`)
		g.P("}")
		g.P()
	}

	g.P("func (Unimplemented", servName, "Server) mustEmbedUnimplemented", servName, "Server() {}")

	g.P()

	g.P("func Register", servName, "Server(r server.ServiceRegistrar, impl ", servName, "Server) {")
	g.P("r.RegisterService(", servName, "_ServiceDesc, impl)")
	g.P("}")

	for _, method := range service.Method {
		g.P("func _", servName, "Server_", method.GetName(), "_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor server.Interceptor) (interface{}, error) {")
		g.P("in := new(", method.GetInputType(), ")")
		g.P("if err := dec(in); err != nil {")
		g.P("return nil, err")
		g.P("}")
		g.P()
		g.P("if interceptor == nil {")
		g.P("return srv.(", servName, "Server).", method.GetName(), "(ctx, in)")
		g.P("}")
		g.P()
		g.P("info := &server.InterceptorServerInfo{")
		g.P("Server:     srv,")
		g.P(`Method: "`, method.GetName(), `",`)
		g.P("}")
		g.P("handler := func(ctx context.Context, req interface{}) (interface{}, error) {")
		g.P("return srv.(", servName, "Server).", method.GetName(), "(ctx, req.(*", method.GetInputType(), "))")
		g.P("}")
		g.P("return interceptor(ctx, in, info, handler)")
		g.P("}")
		g.P()
	}

	g.P("var ", servName, "_ServiceDesc = &server.ServiceDesc{")
	g.P("Methods: []server.MethodDesc{")
	for _, method := range service.Method {
		g.P("{")
		g.P(`Name:    "`, method.GetName(), `",`)
		g.P("Handler: _", servName, "Server_", method.GetName(), "_Handler,")
		g.P("},")
	}
	g.P("},")
	g.P("}")

}
