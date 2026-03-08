package tools

import (
	"context"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterInstanceTools(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("launch",
			mcp.WithDescription("Launch a new Multipass VM instance."),
			mcp.WithString("image", mcp.Description("Image to launch (e.g. \"22.04\", \"daily:24.04\"). Defaults to latest LTS.")),
			mcp.WithString("name", mcp.Description("Name for the instance. Auto-generated if omitted.")),
			mcp.WithNumber("cpus", mcp.Description("Number of CPUs to allocate.")),
			mcp.WithString("disk", mcp.Description("Disk size (e.g. \"10G\", \"50G\").")),
			mcp.WithString("memory", mcp.Description("Memory size (e.g. \"1G\", \"4G\").")),
			mcp.WithString("cloud_init", mcp.Description("Path or URL to cloud-init config file.")),
			mcp.WithString("network", mcp.Description("Network to connect to (from multipass networks).")),
			mcp.WithString("mount", mcp.Description("Host path to mount in the instance (source:target format).")),
			mcp.WithNumber("timeout", mcp.Description("Timeout in seconds (default 600 — image downloads can be slow).")),
		),
		handleLaunch,
	)

	s.AddTool(
		mcp.NewTool("start",
			mcp.WithDescription("Start a stopped Multipass instance."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Instance name, or \"all\" to start all instances.")),
		),
		handleStart,
	)

	s.AddTool(
		mcp.NewTool("stop",
			mcp.WithDescription("Stop a running Multipass instance."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Instance name, or \"all\" to stop all instances.")),
			mcp.WithBoolean("force", mcp.Description("Force stop without waiting for graceful shutdown.")),
		),
		handleStop,
	)

	s.AddTool(
		mcp.NewTool("restart",
			mcp.WithDescription("Restart a Multipass instance."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Instance name, or \"all\" to restart all instances.")),
		),
		handleRestart,
	)

	s.AddTool(
		mcp.NewTool("suspend",
			mcp.WithDescription("Suspend a running Multipass instance."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Instance name, or \"all\" to suspend all instances.")),
		),
		handleSuspend,
	)

	s.AddTool(
		mcp.NewTool("delete",
			mcp.WithDescription("Delete a Multipass instance (can be recovered unless purged)."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Instance name, or \"all\" to delete all instances.")),
			mcp.WithBoolean("purge", mcp.Description("Permanently delete instead of moving to trash.")),
		),
		handleDelete,
	)

	s.AddTool(
		mcp.NewTool("recover",
			mcp.WithDescription("Recover a deleted (trashed) Multipass instance."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Instance name, or \"all\" to recover all deleted instances.")),
		),
		handleRecover,
	)
}

func handleLaunch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := []string{"launch"}

	image := req.GetString("image", "")
	if image != "" {
		args = append(args, image)
	}
	if name := req.GetString("name", ""); name != "" {
		args = append(args, "--name", name)
	}
	if cpus := req.GetInt("cpus", 0); cpus > 0 {
		args = append(args, "--cpus", strconv.Itoa(cpus))
	}
	if disk := req.GetString("disk", ""); disk != "" {
		args = append(args, "--disk", disk)
	}
	if memory := req.GetString("memory", ""); memory != "" {
		args = append(args, "--memory", memory)
	}
	if cloudInit := req.GetString("cloud_init", ""); cloudInit != "" {
		args = append(args, "--cloud-init", cloudInit)
	}
	if network := req.GetString("network", ""); network != "" {
		args = append(args, "--network", network)
	}
	if mount := req.GetString("mount", ""); mount != "" {
		args = append(args, "--mount", mount)
	}

	timeout := launchTimeout
	if t := req.GetInt("timeout", 0); t > 0 {
		timeout = time.Duration(t) * time.Second
	}

	result, err := runMultipass(ctx, timeout, args...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}

func handleStart(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var args []string
	if name == "all" {
		args = []string{"start", "--all"}
	} else {
		args = []string{"start", name}
	}

	result, err := runMultipass(ctx, defaultTimeout, args...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}

func handleStop(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	args := []string{"stop"}
	if req.GetBool("force", false) {
		args = append(args, "--force")
	}
	if name == "all" {
		args = append(args, "--all")
	} else {
		args = append(args, name)
	}

	result, err := runMultipass(ctx, defaultTimeout, args...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}

func handleRestart(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var args []string
	if name == "all" {
		args = []string{"restart", "--all"}
	} else {
		args = []string{"restart", name}
	}

	result, err := runMultipass(ctx, defaultTimeout, args...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}

func handleSuspend(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var args []string
	if name == "all" {
		args = []string{"suspend", "--all"}
	} else {
		args = []string{"suspend", name}
	}

	result, err := runMultipass(ctx, defaultTimeout, args...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}

func handleDelete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	args := []string{"delete"}
	if req.GetBool("purge", false) {
		args = append(args, "--purge")
	}
	if name == "all" {
		args = append(args, "--all")
	} else {
		args = append(args, name)
	}

	result, err := runMultipass(ctx, defaultTimeout, args...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}

func handleRecover(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var args []string
	if name == "all" {
		args = []string{"recover", "--all"}
	} else {
		args = []string{"recover", name}
	}

	result, err := runMultipass(ctx, defaultTimeout, args...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}
