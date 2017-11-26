package funcs

import (
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"github.com/nestgo/nest/util"
)

//FuncMap func map for tempalte render
var FuncMap = map[string]interface{}{
	"LowerFirst": util.LowerFirst,
	"CamelCase":  generator.CamelCase,
}
