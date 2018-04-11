package main

//go:generate go run script/inline_schemas/inline_schemas.go

import (
	"os"

	"github.com/elastic/apm-server/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
