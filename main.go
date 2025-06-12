package main

import (
	_ "embed"
	"os"
	"strings"

	"github.com/netbirdio/management-refactor/cmd"
)

//go:embed go.mod
var GoMod string

func ModulePath() string {
	lines := strings.Split(GoMod, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return ""
}

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
