package main

import (
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/rootisgod/multipass-mcp/tools"
)

func main() {
	s := server.NewMCPServer(
		"multipass",
		"0.2.0",
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

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
