package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterSnapshotTools(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("multipass_snapshot",
			mcp.WithDescription("Create a snapshot of a stopped Multipass instance.\n\nThe instance must be stopped before taking a snapshot."),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Instance name.")),
			mcp.WithString("name", mcp.Description("Snapshot name. Auto-generated if omitted.")),
			mcp.WithString("comment", mcp.Description("Description or comment for the snapshot.")),
			mcp.WithTitleAnnotation("Multipass: Create Snapshot"),
			mcp.WithReadOnlyHintAnnotation(false),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithIdempotentHintAnnotation(false),
			mcp.WithOpenWorldHintAnnotation(false),
		),
		handleSnapshot,
	)

	s.AddTool(
		mcp.NewTool("multipass_restore",
			mcp.WithDescription("Restore a Multipass instance to a snapshot.\n\nThe instance must be stopped. Use destructive=true to discard current state."),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Instance name.")),
			mcp.WithString("snapshot", mcp.Required(), mcp.Description("Snapshot name to restore to.")),
			mcp.WithBoolean("destructive", mcp.Description("Discard current instance state (required if state changed since snapshot).")),
			mcp.WithTitleAnnotation("Multipass: Restore Snapshot"),
			mcp.WithReadOnlyHintAnnotation(false),
			mcp.WithDestructiveHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(false),
		),
		handleRestore,
	)

	s.AddTool(
		mcp.NewTool("multipass_clone",
			mcp.WithDescription("Clone a stopped Multipass instance into an independent copy.\n\nThe source instance must be stopped."),
			mcp.WithString("source_name", mcp.Required(), mcp.Description("Name of the instance to clone.")),
			mcp.WithString("name", mcp.Description("Name for the new cloned instance. Auto-generated if omitted.")),
			mcp.WithTitleAnnotation("Multipass: Clone Instance"),
			mcp.WithReadOnlyHintAnnotation(false),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithIdempotentHintAnnotation(false),
			mcp.WithOpenWorldHintAnnotation(false),
		),
		handleClone,
	)
}

func handleSnapshot(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instance, err := req.RequireString("instance")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	args := []string{"snapshot", instance}
	if name := req.GetString("name", ""); name != "" {
		args = append(args, "--name", name)
	}
	if comment := req.GetString("comment", ""); comment != "" {
		args = append(args, "--comment", comment)
	}

	result, err := runMultipass(ctx, defaultTimeout, args...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}

func handleRestore(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instance, err := req.RequireString("instance")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	snapshot, err := req.RequireString("snapshot")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	args := []string{"restore", fmt.Sprintf("%s.%s", instance, snapshot)}
	if req.GetBool("destructive", false) {
		args = append(args, "--destructive")
	}

	result, err := runMultipass(ctx, defaultTimeout, args...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}

func handleClone(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sourceName, err := req.RequireString("source_name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	args := []string{"clone", sourceName}
	if name := req.GetString("name", ""); name != "" {
		args = append(args, "--name", name)
	}

	result, err := runMultipass(ctx, defaultTimeout, args...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}
