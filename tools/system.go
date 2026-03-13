package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterSystemTools(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("multipass_purge",
			mcp.WithDescription("Permanently delete all trashed (deleted) Multipass instances.\n\nThis is irreversible. All instances previously deleted with 'multipass_delete' will be permanently removed."),
			mcp.WithTitleAnnotation("Multipass: Purge Deleted Instances"),
			mcp.WithReadOnlyHintAnnotation(false),
			mcp.WithDestructiveHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(false),
		),
		handlePurge,
	)

	s.AddTool(
		mcp.NewTool("multipass_authenticate",
			mcp.WithDescription("Authenticate with the Multipass service using a passphrase."),
			mcp.WithString("passphrase", mcp.Required(), mcp.Description("Authentication passphrase.")),
			mcp.WithTitleAnnotation("Multipass: Authenticate"),
			mcp.WithReadOnlyHintAnnotation(false),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithIdempotentHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(false),
		),
		handleAuthenticate,
	)
}

func handlePurge(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result, err := runMultipass(ctx, defaultTimeout, "purge")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}

func handleAuthenticate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	passphrase, err := req.RequireString("passphrase")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, err := runMultipass(ctx, defaultTimeout, "authenticate", passphrase)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}
