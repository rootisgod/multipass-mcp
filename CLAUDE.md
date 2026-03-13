# CLAUDE.md

## Project Overview

Multipass MCP Server — exposes the Multipass CLI as an MCP server with 32 tools and 6 resources. Lets AI assistants manage Ubuntu VMs.

## Tech Stack

- Go 1.23+, static binary (CGO_ENABLED=0)
- mark3labs/mcp-go for MCP server framework
- os/exec for subprocess calls to `multipass` CLI
- GoReleaser for cross-platform builds and releases

## Project Structure

```
main.go                     # entrypoint, server setup, tool registration, version via ldflags
tools/
├── helpers.go              # runMultipass/runMultipassJSON subprocess helpers
├── instance.go             # multipass_launch, multipass_start, multipass_stop, multipass_restart, multipass_suspend, multipass_delete, multipass_recover
├── exec.go                 # multipass_run_script, multipass_exec_command
├── files.go                # multipass_transfer, multipass_mount_directory, multipass_umount_directory
├── snapshots.go            # multipass_snapshot, multipass_restore, multipass_clone
├── config.go               # multipass_get_config, multipass_set_config
├── system.go               # multipass_purge, multipass_authenticate
├── info.go                 # multipass_list_instances, multipass_get_instance, multipass_find_images, multipass_list_networks, multipass_list_snapshots, multipass_get_version, multipass_list_aliases, multipass_list_mounts, multipass_list_deleted, multipass_instance_exists, multipass_get_bridged_network, multipass_disk_usage_check, multipass_wait_until_running
└── resources.go            # 6 resources (5 static + 1 template)
go.mod / go.sum             # module dependencies
Dockerfile                  # multi-stage Go build
.goreleaser.yaml            # cross-compile + Cosign signing + Homebrew tap
.github/workflows/release.yml  # GitHub Actions release pipeline
```

## Key Patterns

- `runMultipass(ctx, timeout, args...)` — subprocess helper, returns stdout or error with stderr
- `runMultipassJSON(ctx, timeout, args...)` — same but appends `--format json` and parses result
- Resources registered via `s.AddResource()` / `s.AddResourceTemplate()`
- Tools registered via `s.AddTool(tool, handler)` with `multipass_` prefix
- All tools include MCP annotations (readOnlyHint, destructiveHint, idempotentHint, openWorldHint)
- `name="all"` on lifecycle tools translates to `--all` CLI flag
- `multipass_exec_command` takes `command: []string` (not a shell string) to prevent injection
- Timeouts: 300s default, 600s for `multipass_launch`
- Version injected at build time via `-ldflags "-X main.version=..."`

## Development Commands

```bash
# Build
CGO_ENABLED=0 go build -ldflags "-X main.version=$(git describe --tags --always)" -o multipass-mcp .

# Check
go vet ./...

# Register with Claude Code
claude mcp add multipass-mcp -- /path/to/multipass-mcp
```

## Entry Point

`main.go` → registers tools/resources on `server.MCPServer` → `server.ServeStdio(s)`

## Verified Working

End-to-end tested: launch VM, install packages via multipass_exec_command, deploy files, delete+purge. All 32 tools and 6 resources confirmed registered (5 static + 1 template).

## Notes

- `multipass_exec_command` can run complex shell commands via `["bash", "-c", "..."]`
- CLI spinner/escape codes appear in tool output but don't affect functionality
- Empty string result from tools (e.g. delete, mount) means success
- Ubuntu 24.04 (Noble) is the current default LTS image on arm64
