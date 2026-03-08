package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

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

	s.AddTool(
		mcp.NewTool("list_mounts",
			mcp.WithDescription("List active mounts for a Multipass instance. Returns mount source, target, UID/GID mappings."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Instance name.")),
		),
		handleListMountsTool,
	)

	s.AddTool(
		mcp.NewTool("list_deleted",
			mcp.WithDescription("List only deleted (trashed) Multipass instances that can be recovered or purged."),
		),
		handleListDeletedTool,
	)

	s.AddTool(
		mcp.NewTool("instance_exists",
			mcp.WithDescription("Check if a Multipass instance exists by name. Returns exists (bool), state, and image. Use before launch to avoid duplicates or before lifecycle commands to verify the target."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Instance name to check.")),
		),
		handleInstanceExistsTool,
	)

	s.AddTool(
		mcp.NewTool("get_bridged_network",
			mcp.WithDescription("Get the configured bridged network interface name. Returns the network used when launching with --network bridged."),
		),
		handleGetBridgedNetworkTool,
	)

	s.AddTool(
		mcp.NewTool("disk_usage_check",
			mcp.WithDescription("Check disk usage for a Multipass instance. Returns total, used, percentage, and a warning if usage exceeds the threshold."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Instance name.")),
			mcp.WithNumber("warn_percent", mcp.Description("Warning threshold percentage (default 80).")),
		),
		handleDiskUsageCheckTool,
	)

	s.AddTool(
		mcp.NewTool("wait_until_running",
			mcp.WithDescription("Wait for a Multipass instance to reach Running state. Polls until the instance is running or timeout is reached. Useful after launch, start, or restart."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Instance name.")),
			mcp.WithNumber("timeout", mcp.Description("Max seconds to wait (default 120).")),
		),
		handleWaitUntilRunningTool,
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

func handleListMountsTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	data, err := runMultipassJSON(ctx, defaultTimeout, "info", name)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var info map[string]any
	if err := json.Unmarshal(data, &info); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse info: %v", err)), nil
	}

	// Extract mounts from the info response
	mounts := make(map[string]any)
	if infoMap, ok := info["info"].(map[string]any); ok {
		if inst, ok := infoMap[name].(map[string]any); ok {
			if m, ok := inst["mounts"].(map[string]any); ok {
				mounts = m
			}
		}
	}

	pretty, _ := json.MarshalIndent(mounts, "", "  ")
	return mcp.NewToolResultText(string(pretty)), nil
}

func handleListDeletedTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	data, err := runMultipassJSON(ctx, defaultTimeout, "list")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var list map[string]any
	if err := json.Unmarshal(data, &list); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse list: %v", err)), nil
	}

	var deleted []any
	if instances, ok := list["list"].([]any); ok {
		for _, inst := range instances {
			if m, ok := inst.(map[string]any); ok {
				if state, ok := m["state"].(string); ok && state == "Deleted" {
					deleted = append(deleted, m)
				}
			}
		}
	}

	result := map[string]any{"deleted": deleted}
	pretty, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(pretty)), nil
}

func handleInstanceExistsTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	data, err := runMultipassJSON(ctx, defaultTimeout, "list")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var list map[string]any
	if err := json.Unmarshal(data, &list); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse list: %v", err)), nil
	}

	result := map[string]any{"name": name, "exists": false}
	if instances, ok := list["list"].([]any); ok {
		for _, inst := range instances {
			if m, ok := inst.(map[string]any); ok {
				if n, ok := m["name"].(string); ok && n == name {
					result["exists"] = true
					if state, ok := m["state"].(string); ok {
						result["state"] = state
					}
					if release, ok := m["release"].(string); ok {
						result["image"] = release
					}
					break
				}
			}
		}
	}

	pretty, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(pretty)), nil
}

func handleGetBridgedNetworkTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result, err := runMultipass(ctx, defaultTimeout, "get", "local.bridged-network")
	if err != nil {
		if strings.Contains(err.Error(), "is not a valid key") || strings.Contains(err.Error(), "not set") {
			return mcp.NewToolResultText(`{"bridged_network": null, "message": "No bridged network configured. Use set_config to set local.bridged-network."}`), nil
		}
		return mcp.NewToolResultError(err.Error()), nil
	}
	resp := map[string]string{"bridged_network": strings.TrimSpace(result)}
	pretty, _ := json.MarshalIndent(resp, "", "  ")
	return mcp.NewToolResultText(string(pretty)), nil
}

func handleDiskUsageCheckTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	warnPercent := 80.0
	if w := req.GetInt("warn_percent", 0); w > 0 {
		warnPercent = float64(w)
	}

	data, err := runMultipassJSON(ctx, defaultTimeout, "info", name)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var info map[string]any
	if err := json.Unmarshal(data, &info); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse info: %v", err)), nil
	}

	result := map[string]any{"name": name}
	if infoMap, ok := info["info"].(map[string]any); ok {
		if inst, ok := infoMap[name].(map[string]any); ok {
			if disks, ok := inst["disks"].(map[string]any); ok {
				if sda, ok := disks["sda1"].(map[string]any); ok {
					total, _ := sda["total"].(string)
					used, _ := sda["used"].(string)
					result["total"] = total
					result["used"] = used

					// Parse byte values for percentage calculation
					totalBytes := parseDiskBytes(total)
					usedBytes := parseDiskBytes(used)
					if totalBytes > 0 {
						pct := (float64(usedBytes) / float64(totalBytes)) * 100
						result["percent_used"] = fmt.Sprintf("%.1f%%", pct)
						result["warning"] = pct >= warnPercent
						if pct >= warnPercent {
							result["message"] = fmt.Sprintf("Disk usage %.1f%% exceeds threshold %.0f%%", pct, warnPercent)
						}
					}
				}
			}
		}
	}

	pretty, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(pretty)), nil
}

// parseDiskBytes parses multipass disk size strings like "5765507072" to int64.
func parseDiskBytes(s string) int64 {
	s = strings.TrimSpace(s)
	var val int64
	fmt.Sscanf(s, "%d", &val)
	return val
}

func handleWaitUntilRunningTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	timeoutSecs := 120
	if t := req.GetInt("timeout", 0); t > 0 {
		timeoutSecs = t
	}

	deadline := time.Now().Add(time.Duration(timeoutSecs) * time.Second)
	pollInterval := 3 * time.Second

	for {
		data, err := runMultipassJSON(ctx, defaultTimeout, "info", name)
		if err != nil {
			if time.Now().After(deadline) {
				return mcp.NewToolResultError(fmt.Sprintf("timeout after %ds waiting for %s: %v", timeoutSecs, name, err)), nil
			}
			time.Sleep(pollInterval)
			continue
		}

		var info map[string]any
		if err := json.Unmarshal(data, &info); err == nil {
			if infoMap, ok := info["info"].(map[string]any); ok {
				if inst, ok := infoMap[name].(map[string]any); ok {
					if state, ok := inst["state"].(string); ok {
						if state == "Running" {
							ipv4 := ""
							if addrs, ok := inst["ipv4"].([]any); ok && len(addrs) > 0 {
								if ip, ok := addrs[0].(string); ok {
									ipv4 = ip
								}
							}
							result := map[string]any{
								"name":  name,
								"state": "Running",
								"ipv4":  ipv4,
							}
							pretty, _ := json.MarshalIndent(result, "", "  ")
							return mcp.NewToolResultText(string(pretty)), nil
						}
						if state == "Deleted" || state == "Unknown" {
							return mcp.NewToolResultError(fmt.Sprintf("instance %s is in %s state and will not reach Running", name, state)), nil
						}
					}
				}
			}
		}

		if time.Now().After(deadline) {
			return mcp.NewToolResultError(fmt.Sprintf("timeout after %ds waiting for %s to reach Running state", timeoutSecs, name)), nil
		}
		time.Sleep(pollInterval)
	}
}
