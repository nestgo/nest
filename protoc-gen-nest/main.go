package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
)

func main() {
	// Begin by allocating a generator. The request and response structures are stored there
	// so we can do error handling easily - the response structure contains the field to
	// report failure.
	gen := generator.New()

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		gen.Error(err, "reading input")
	}

	if err := proto.Unmarshal(data, gen.Request); err != nil {
		gen.Error(err, "parsing input proto")
	}

	if len(gen.Request.FileToGenerate) == 0 {
		gen.Fail("no files to generate")
	}

	gen.CommandLineParameters(gen.Request.GetParameter())
	// Create a wrapped version of the Descriptors and EnumDescriptors that
	// point to the file that defines them.
	gen.WrapTypes()

	gen.SetPackageNames()
	gen.BuildTypeNameMap()
	gen.GenerateAllFiles()
	gen.Response.File = []*plugin_go.CodeGeneratorResponse_File{}

	g := new(grpc)
	g.gen = gen

	var destinationDir = "."
	//parse params
	if parameter := gen.Request.GetParameter(); parameter != "" {
		for _, param := range strings.Split(parameter, ",") {
			parts := strings.Split(param, "=")
			if len(parts) != 2 {
				log.Printf("Err: invalid parameter: %q", param)
				continue
			}
			switch parts[0] {
			case "output":
				destinationDir = parts[1]
				break
			default:
				log.Printf("Err: unknown parameter: %q", param)
			}
		}
	}

	tmplMap := make(map[string]*plugin_go.CodeGeneratorResponse_File)
	concatOrAppend := func(file *plugin_go.CodeGeneratorResponse_File) {
		if val, ok := tmplMap[file.GetName()]; ok {
			*val.Content += file.GetContent()
		} else {
			tmplMap[file.GetName()] = file
			gen.Response.File = append(gen.Response.File, file)
		}
	}

	// Generate the encoders
	for _, file := range gen.Request.GetProtoFile() {
		encoder := NewGenericTemplateBasedEncoder(file, destinationDir, g)
		for _, tmpl := range encoder.Files() {
			concatOrAppend(tmpl)
		}
	}

	// Send back the results.
	data, err = proto.Marshal(gen.Response)
	if err != nil {
		gen.Error(err, "failed to marshal output proto")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		gen.Error(err, "failed to write output proto")
	}
}
