package main

import (
	"../config"
	"../server"
	"github.com/nestgo/nest"
)

func main() {
	cfg := config.GetConfig()

	nest.Register(
		{{- range .File.Service}}
		&server.{{CamelCase .Name}}Server{},
		{{end}}
	)
	nest.Run(cfg.Address)
}
