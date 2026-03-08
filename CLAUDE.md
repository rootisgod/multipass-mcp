# CLAUDE.md

## Project Overview

Multipass MCP Server — exposes the Multipass CLI as an MCP server with 25 tools and 6 resources. Lets AI assistants manage Ubuntu VMs.

## Tech Stack

- Go 1.23+, static binary (CGO_ENABLED=0)
- mark3labs/mcp-go for MCP server framework
- os/exec for subprocess calls to `multipass` CLI
- GoReleaser for cross-platform builds and releases

## Project Structure

```
main.go                     # entrypoint, server setup, tool registration
tools/
├── helpers.go              # runMultipass/runMultipassJSON subprocess helpers
├── instance.go             # launch, start, stop, restart, suspend, delete, recover
├── exec.go                 # exec_command
├── files.go                # transfer, mount_directory, umount_directory
├── snapshots.go            # snapshot, restore, clone
├── config.go               # get_config, set_config
├── system.go               # purge, authenticate
├── info.go                 # list_instances, get_instance, find_images, list_networks, list_snapshots, get_version, list_aliases
└── resources.go            # 6 resources (5 static + 1 template)
go.mod / go.sum             # module dependencies
.goreleaser.yaml            # cross-compile + Cosign signing + Homebrew tap
.github/workflows/release.yml  # GitHub Actions release pipeline
```

## Key Patterns

- `runMultipass(ctx, timeout, args...)` — subprocess helper, returns stdout or error with stderr
- `runMultipassJSON(ctx, timeout, args...)` — same but appends `--format json` and parses result
- Resources registered via `s.AddResource()` / `s.AddResourceTemplate()`
- Tools registered via `s.AddTool(tool, handler)`
- `name="all"` on lifecycle tools translates to `--all` CLI flag
- `exec_command` takes `command: []string` (not a shell string) to prevent injection
- Timeouts: 300s default, 600s for `launch`

## Development Commands

```bash
# Build
CGO_ENABLED=0 go build -o multipass-mcp .

# Check
go vet ./...

# Register with Claude Code
claude mcp add multipass-mcp -- /path/to/multipass-mcp
```

## Entry Point

`main.go` → registers tools/resources on `server.MCPServer` → `server.ServeStdio(s)`

## Verified Working

End-to-end tested: launch VM, install packages via exec_command, deploy files, delete+purge. All 25 tools and 6 resources confirmed registered (5 static + 1 template).

## Notes

- `exec_command` can run complex shell commands via `["bash", "-c", "..."]`
- CLI spinner/escape codes appear in tool output but don't affect functionality
- Empty string result from tools (e.g. delete, mount) means success
- Ubuntu 24.04 (Noble) is the current default LTS image on arm64
