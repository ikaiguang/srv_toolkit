package main

import (
	"flag"
	"github.com/ikaiguang/srv_toolkit/cmd/protoc-gen-go-tkgrpc/gengogrpc"

	gengo "google.golang.org/protobuf/cmd/protoc-gen-go/internal_gengo"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	var (
		flags        flag.FlagSet
		//plugins      = flags.String("plugins", "", "list of plugins to enable (supported values: grpc)")
		importPrefix = flags.String("import_prefix", "", "prefix to prepend to import paths")
	)
	importRewriteFunc := func(importPath protogen.GoImportPath) protogen.GoImportPath {
		switch importPath {
		case "context", "fmt", "math":
			return importPath
		}
		if *importPrefix != "" {
			return protogen.GoImportPath(*importPrefix) + importPath
		}
		return importPath
	}
	protogen.Options{
		ParamFunc:         flags.Set,
		ImportRewriteFunc: importRewriteFunc,
	}.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			gengogrpc.GenerateFile(gen, f)
		}
		gen.SupportedFeatures = gengo.SupportedFeatures
		return nil
	})
}