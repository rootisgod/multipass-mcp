package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterExecTools(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("exec_command",
			mcp.WithDescription("Execute a command inside a Multipass instance."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Instance name.")),
			mcp.WithArray("command",
				mcp.Required(),
				mcp.Description("Command and arguments as a list (e.g. [\"ls\", \"-la\", \"/tmp\"])."),
				mcp.Items(map[string]any{"type": "string"}),
			),
			mcp.WithString("working_directory", mcp.Description("Working directory inside the instance.")),
		),
		handleExecCommand,
	)
}

func handleExecCommand(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	rawArgs := req.GetArguments()
	commandRaw, ok := rawArgs["command"]
	if !ok {
		return mcp.NewToolResultError("missing required parameter: command"), nil
	}

	commandSlice, ok := commandRaw.([]any)
	if !ok {
		return mcp.NewToolResultError("command must be a list of strings"), nil
	}

	cmdArgs := []string{"exec", name}

	if wd := req.GetString("working_directory", ""); wd != "" {
		cmdArgs = append(cmdArgs, "--working-directory", wd)
	}

	cmdArgs = append(cmdArgs, "--")
	for _, item := range commandSlice {
		s, ok := item.(string)
		if !ok {
			return mcp.NewToolResultError("each command element must be a string"), nil
		}
		cmdArgs = append(cmdArgs, s)
	}

	result, err := runMultipass(ctx, defaultTimeout, cmdArgs...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}
