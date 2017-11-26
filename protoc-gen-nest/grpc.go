package main

import "github.com/golang/protobuf/protoc-gen-go/generator"

// grpc is an implementation of the Go protocol buffer compiler's
// plugin architecture.  It generates bindings for gRPC support.
type grpc struct {
	gen *generator.Generator
}

// Init initializes the plugin.
func (g *grpc) Init(gen *generator.Generator) {
	g.gen = gen
}

// Given a type name defined in a .proto, return its object.
// Also record that we're using it, to guarantee the associated import.
func (g *grpc) objectNamed(name string) generator.Object {
	g.gen.RecordTypeUse(name)
	return g.gen.ObjectNamed(name)
}

// Given a type name defined in a .proto, return its name as we will print it.
func (g *grpc) typeName(str string) string {
	return g.gen.TypeName(g.objectNamed(str))
}
