package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterFileTools(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("multipass_transfer",
			mcp.WithDescription("Transfer files between host and a Multipass instance.\n\nUse <instance>:<path> syntax for instance paths (e.g. \"my-vm:/home/ubuntu/file.txt\").\nUse a plain path for host paths (e.g. \"/tmp/file.txt\")."),
			mcp.WithString("source", mcp.Required(), mcp.Description("Source path (host path or instance:path).")),
			mcp.WithString("destination", mcp.Required(), mcp.Description("Destination path (host path or instance:path).")),
			mcp.WithBoolean("recursive", mcp.Description("Transfer directories recursively.")),
			mcp.WithTitleAnnotation("Multipass: Transfer Files"),
			mcp.WithReadOnlyHintAnnotation(false),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithIdempotentHintAnnotation(false),
			mcp.WithOpenWorldHintAnnotation(false),
		),
		handleTransfer,
	)

	s.AddTool(
		mcp.NewTool("multipass_mount_directory",
			mcp.WithDescription("Mount a host directory inside a Multipass instance."),
			mcp.WithString("source", mcp.Required(), mcp.Description("Host directory path to mount.")),
			mcp.WithString("target", mcp.Required(), mcp.Description("Mount point in instance:path format (e.g. \"my-vm:/mnt/data\").")),
			mcp.WithString("uid_map", mcp.Description("UID mapping in host:instance format (e.g. \"1000:0\").")),
			mcp.WithString("gid_map", mcp.Description("GID mapping in host:instance format (e.g. \"1000:0\").")),
			mcp.WithString("mount_type", mcp.Description("Mount type: \"classic\" or \"native\".")),
			mcp.WithTitleAnnotation("Multipass: Mount Directory"),
			mcp.WithReadOnlyHintAnnotation(false),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithIdempotentHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(false),
		),
		handleMountDirectory,
	)

	s.AddTool(
		mcp.NewTool("multipass_umount_directory",
			mcp.WithDescription("Unmount a directory from a Multipass instance."),
			mcp.WithString("mount_path", mcp.Required(), mcp.Description("Mount point to remove (instance:path format or host path).")),
			mcp.WithTitleAnnotation("Multipass: Unmount Directory"),
			mcp.WithReadOnlyHintAnnotation(false),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithIdempotentHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(false),
		),
		handleUmountDirectory,
	)
}

func handleTransfer(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	source, err := req.RequireString("source")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	destination, err := req.RequireString("destination")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	args := []string{"transfer"}
	if req.GetBool("recursive", false) {
		args = append(args, "--recursive")
	}
	args = append(args, source, destination)

	result, err := runMultipass(ctx, defaultTimeout, args...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}

func handleMountDirectory(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	source, err := req.RequireString("source")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	target, err := req.RequireString("target")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	args := []string{"mount", source, target}

	if uidMap := req.GetString("uid_map", ""); uidMap != "" {
		args = append(args, "--uid-map", uidMap)
	}
	if gidMap := req.GetString("gid_map", ""); gidMap != "" {
		args = append(args, "--gid-map", gidMap)
	}
	if mountType := req.GetString("mount_type", ""); mountType != "" {
		args = append(args, "--type", mountType)
	}

	result, err := runMultipass(ctx, defaultTimeout, args...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}

func handleUmountDirectory(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	mountPath, err := req.RequireString("mount_path")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, err := runMultipass(ctx, defaultTimeout, "umount", mountPath)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}
