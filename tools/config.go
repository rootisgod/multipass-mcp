package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterConfigTools(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("get_config",
			mcp.WithDescription("Get Multipass configuration settings."),
			mcp.WithString("key", mcp.Description("Specific setting key (e.g. \"local.driver\"). Omit to get all settings.")),
		),
		handleGetConfig,
	)

	s.AddTool(
		mcp.NewTool("set_config",
			mcp.WithDescription("Set a Multipass configuration value."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Setting key (e.g. \"local.driver\").")),
			mcp.WithString("value", mcp.Required(), mcp.Description("New value to set.")),
		),
		handleSetConfig,
	)
}

func handleGetConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key := req.GetString("key", "")

	var args []string
	if key != "" {
		args = []string{"get", key}
	} else {
		args = []string{"get", "--keys"}
	}

	result, err := runMultipass(ctx, defaultTimeout, args...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}

func handleSetConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	value, err := req.RequireString("value")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, err := runMultipass(ctx, defaultTimeout, "set", key+"="+value)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}
