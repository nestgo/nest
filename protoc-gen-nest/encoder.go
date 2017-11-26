package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/nestgo/nest/protoc-gen-nest/funcs"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/nestgo/nest/util"
)

//GenericTemplateBasedEncoder struct
type GenericTemplateBasedEncoder struct {
	g              *grpc
	service        *descriptor.ServiceDescriptorProto
	file           *descriptor.FileDescriptorProto
	enum           []*descriptor.EnumDescriptorProto
	destinationDir string
}

//Ast ast
type Ast struct {
	BuildDate        time.Time                          `json:"build-date"`
	BuildHostname    string                             `json:"build-hostname"`
	BuildUser        string                             `json:"build-user"`
	GoPWD            string                             `json:"go-pwd,omitempty"`
	PWD              string                             `json:"pwd"`
	DestinationDir   string                             `json:"destination-dir"`
	File             *descriptor.FileDescriptorProto    `json:"file"`
	RawFilename      string                             `json:"raw-filename"`
	Filename         string                             `json:"filename"`
	TemplateDir      string                             `json:"template-dir"`
	Service          *descriptor.ServiceDescriptorProto `json:"service"`
	Enum             []*descriptor.EnumDescriptorProto  `json:"enum"`
	PackageName      string                             `json:"package-name"`
	ClientMethodsTpl string                             `json:"client-methods-tpl"`
	ServerMethodsTpl string                             `json:"server-methods-tpl"`
}

//NewGenericTemplateBasedEncoder new
func NewGenericTemplateBasedEncoder(file *descriptor.FileDescriptorProto, destinationDir string, g *grpc) (e *GenericTemplateBasedEncoder) {
	e = &GenericTemplateBasedEncoder{
		g:              g,
		service:        nil,
		file:           file,
		enum:           file.GetEnumType(),
		destinationDir: destinationDir,
	}
	return
}

func (e *GenericTemplateBasedEncoder) genAst(templateFilename string, service *descriptor.ServiceDescriptorProto) (*Ast, error) {
	// prepare the ast passed to the template engine
	hostname, _ := os.Hostname()
	pwd, _ := os.Getwd()
	goPwd := ""
	if os.Getenv("GOPATH") != "" {
		goPwd, _ = filepath.Rel(os.Getenv("GOPATH")+"/src", pwd)
		if strings.Contains(goPwd, "../") {
			goPwd = ""
		}
	}
	packageName, err := e.getPackageName()
	if err != nil {
		return nil, err
	}
	ast := &Ast{
		BuildDate:      time.Now(),
		BuildHostname:  hostname,
		BuildUser:      os.Getenv("USER"),
		PWD:            pwd,
		GoPWD:          goPwd,
		File:           e.file,
		DestinationDir: e.destinationDir,
		RawFilename:    templateFilename,
		Service:        service,
		Enum:           e.enum,
		PackageName:    packageName,
	}
	if service != nil {
		ast.ClientMethodsTpl = e.GenClientMethodsCode(service)
		ast.ServerMethodsTpl = e.GenServerMethodsCode(service)
	}
	return ast, nil
}

func (e *GenericTemplateBasedEncoder) buildContent(templateFilename string, service *descriptor.ServiceDescriptorProto) (string, error) {
	fileContentBytes, err := Asset(templateFilename)
	if err != nil {
		return "", err
	}
	// initialize template engine
	tmpl, err := template.New(templateFilename).Funcs(funcs.FuncMap).Parse(string(fileContentBytes))
	if err != nil {
		return "", err
	}

	ast, err := e.genAst(templateFilename, service)
	if err != nil {
		return "", err
	}

	// generate the content
	buffer := new(bytes.Buffer)
	if err := tmpl.Execute(buffer, ast); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

//Files all tempalte file
func (e *GenericTemplateBasedEncoder) Files() []*plugin_go.CodeGeneratorResponse_File {
	templates := AssetNames()

	length := len(templates) + len(e.file.Service)
	files := make([]*plugin_go.CodeGeneratorResponse_File, 0, length)
	for _, tmplFile := range templates {
		if tmplFile == "client.go.tmpl" || tmplFile == "server/server.go.tmpl" {
			for _, service := range e.file.Service {
				copyTmplFile := tmplFile
				switch tmplFile {
				case "client.go.tmpl":
					copyTmplFile = fmt.Sprintf("%s.go", util.LowerFirst(service.GetName()))
				case "server/server.go.tmpl":
					copyTmplFile = fmt.Sprintf("server/%s.go", util.LowerFirst(service.GetName()))
				}
				content, err := e.buildContent(tmplFile, service)
				if err != nil {
					panic(err)
				}
				copyTmplFile = strings.TrimSuffix(copyTmplFile, ".tmpl")
				files = append(files, &plugin_go.CodeGeneratorResponse_File{
					Content: &content,
					Name:    &copyTmplFile,
				})
			}
		} else {
			content, err := e.buildContent(tmplFile, nil)
			if err != nil {
				panic(err)
			}
			copyTmplFile := tmplFile
			copyTmplFile = strings.TrimSuffix(copyTmplFile, ".tmpl")
			files = append(files, &plugin_go.CodeGeneratorResponse_File{
				Content: &content,
				Name:    &copyTmplFile,
			})
		}
	}
	return files
}

func (e *GenericTemplateBasedEncoder) getPackageName() (string, error) {
	path, err := filepath.Abs(e.destinationDir)
	if err != nil {
		return "", err
	}
	paths := strings.Split(path, "/")
	return paths[len(paths)-1], nil
}

//GenClientMethodsCode  get client methods code
func (e *GenericTemplateBasedEncoder) GenClientMethodsCode(service *descriptor.ServiceDescriptorProto) string {
	var tpl string
	origServName := service.GetName()
	servName := generator.CamelCase(origServName)

	for _, method := range service.Method {
		var methodTpl string
		origMethName := method.GetName()
		methName := generator.CamelCase(origMethName)
		reqArg := ", in *proto." + e.g.typeName(method.GetInputType())
		callReqArg := ", in"
		if method.GetClientStreaming() {
			reqArg = ""
			callReqArg = ""
		}
		respName := "*proto." + e.g.typeName(method.GetOutputType())
		if method.GetServerStreaming() || method.GetClientStreaming() {
			respName = "proto." + servName + "_" + generator.CamelCase(origMethName) + "Client"
		}
		methodTpl = fmt.Sprintf("//%s call %s\n", methName, methName)
		methodTpl += fmt.Sprintf("func (%sClient *GreeterClient) %s(ctx context.Context%s) (%s, error) {\n", util.LowerFirst(servName), methName, reqArg, respName)
		methodTpl += fmt.Sprintf("    return %sClient.client.%s(ctx%s)\n", util.LowerFirst(servName), methName, callReqArg)
		methodTpl += "}\n"
		tpl += methodTpl
	}
	return tpl
}

//GenServerMethodsCode get server methods code
func (e *GenericTemplateBasedEncoder) GenServerMethodsCode(service *descriptor.ServiceDescriptorProto) string {
	var tpl string
	origServName := service.GetName()
	servName := generator.CamelCase(origServName)

	for _, method := range service.Method {
		var methodTpl string
		origMethName := method.GetName()
		methName := generator.CamelCase(origMethName)

		var reqArgs []string
		ret := "error"
		returnCode := "nil"
		if !method.GetServerStreaming() && !method.GetClientStreaming() {
			reqArgs = append(reqArgs, "ctx context.Context")
			ret = "(*proto." + e.g.typeName(method.GetOutputType()) + ", error)"
			returnCode = "nil,nil"
		}
		if !method.GetClientStreaming() {
			reqArgs = append(reqArgs, "in *proto."+e.g.typeName(method.GetInputType()))
		}
		if method.GetServerStreaming() || method.GetClientStreaming() {
			reqArgs = append(reqArgs, "serv proto."+servName+"_"+generator.CamelCase(origMethName)+"Server")
		}
		methodTpl = "//" + methName + " call " + methName + "\n"
		methodTpl += "func (" + util.LowerFirst(servName) + " *" + generator.CamelCase(servName) + "Server) " + methName + "(" + strings.Join(reqArgs, ", ") + ") " + ret + " {\n"
		methodTpl += "    //TODO \n"
		methodTpl += "    return " + returnCode + "\n"
		methodTpl += "}\n"
		tpl += methodTpl
	}
	return tpl
}
