package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterInfoTools(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("list_instances",
			mcp.WithDescription("List all Multipass instances with their state, IP addresses, and image info."),
		),
		handleListInstancesTool,
	)

	s.AddTool(
		mcp.NewTool("get_instance",
			mcp.WithDescription("Get detailed information about a specific Multipass instance including CPU, memory, disk usage, mounts, and snapshots."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Instance name.")),
		),
		handleGetInstanceTool,
	)

	s.AddTool(
		mcp.NewTool("find_images",
			mcp.WithDescription("List available Multipass images that can be launched (Ubuntu releases, blueprints, etc)."),
		),
		handleFindImagesTool,
	)

	s.AddTool(
		mcp.NewTool("list_networks",
			mcp.WithDescription("List host network interfaces available for Multipass instances."),
		),
		handleListNetworksTool,
	)

	s.AddTool(
		mcp.NewTool("list_snapshots",
			mcp.WithDescription("List all snapshots for a specific Multipass instance."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Instance name.")),
		),
		handleListSnapshotsTool,
	)

	s.AddTool(
		mcp.NewTool("get_version",
			mcp.WithDescription("Get Multipass version and daemon information."),
		),
		handleGetVersionTool,
	)

	s.AddTool(
		mcp.NewTool("list_aliases",
			mcp.WithDescription("List configured Multipass command aliases."),
		),
		handleListAliasesTool,
	)
}

func jsonToolResult(ctx context.Context, args ...string) (*mcp.CallToolResult, error) {
	data, err := runMultipassJSON(ctx, defaultTimeout, args...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	pretty, _ := json.MarshalIndent(json.RawMessage(data), "", "  ")
	return mcp.NewToolResultText(string(pretty)), nil
}

func handleListInstancesTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return jsonToolResult(ctx, "list")
}

func handleGetInstanceTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonToolResult(ctx, "info", name)
}

func handleFindImagesTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return jsonToolResult(ctx, "find")
}

func handleListNetworksTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return jsonToolResult(ctx, "networks")
}

func handleListSnapshotsTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonToolResult(ctx, "info", name, "--snapshots")
}

func handleGetVersionTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return jsonToolResult(ctx, "version")
}

func handleListAliasesTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return jsonToolResult(ctx, "aliases")
}
