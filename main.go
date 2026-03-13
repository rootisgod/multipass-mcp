package main

import (
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/rootisgod/multipass-mcp/tools"
)

// version is set at build time via -ldflags "-X main.version=..."
var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println("multipass-mcp " + version)
		return
	}

	s := server.NewMCPServer(
		"multipass",
		version,
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
	)

	tools.RegisterResources(s)
	tools.RegisterInstanceTools(s)
	tools.RegisterExecTools(s)
	tools.RegisterFileTools(s)
	tools.RegisterSnapshotTools(s)
	tools.RegisterConfigTools(s)
	tools.RegisterSystemTools(s)
	tools.RegisterInfoTools(s)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
